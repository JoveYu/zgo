package httplog

import (
	"net/http"
	"time"

	"github.com/JoveYu/zgo/log"
)

var DefaultLogRequest = func(start time.Time, req *http.Request, resp *http.Response, err error) {

	if err == nil {
		log.Info("ep=http|method=%s|url=%s|code=%d|req=%d|resp=%d|time=%d",
			req.Method, req.URL, resp.StatusCode, req.ContentLength, resp.ContentLength,
			time.Now().Sub(start)/time.Microsecond,
		)
	} else {
		log.Warn("ep=http|method=%s|url=%s|code=%d|req=%d|resp=%d|time=%d|err=%s",
			req.Method, req.URL, 0, req.ContentLength, 0,
			time.Now().Sub(start)/time.Microsecond, err,
		)
	}
}

var DefaultTransport = &Transport{
	RoundTripper: http.DefaultTransport,
}

type Transport struct {
	http.RoundTripper
	LogRequest func(time.Time, *http.Request, *http.Response, error)
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	resp, err := t.RoundTripper.RoundTrip(req)
	if t.LogRequest != nil {
		t.LogRequest(start, req, resp, err)
	} else {
		DefaultLogRequest(start, req, resp, err)
	}

	return resp, err
}
