package http

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/JoveYu/zgo/log"
)

func NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, url, body)
}

type Client struct {
	*http.Client
}

func NewClient() *Client {
	return &Client{
		Client: &http.Client{},
	}
}

func (c *Client) timit(start time.Time, resp *http.Response, err error) {
	req := resp.Request
	if err == nil {
		log.Info("ep=httpclient|method=%s|url=%s|code=%d|req=%d|resp=%d|time=%d",
			req.Method, req.URL, resp.StatusCode, req.ContentLength, resp.ContentLength,
			time.Now().Sub(start)/time.Microsecond,
		)
	} else {
		log.Warn("ep=httpclient|method=%s|url=%s|code=%d|req=%d|resp=%d|time=%d|err=%s",
			req.Method, req.URL, 0, req.ContentLength, 0,
			time.Now().Sub(start)/time.Microsecond, err,
		)
	}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	start := time.Now()
	resp, err := c.Client.Do(req)
	c.timit(start, resp, err)
	return resp, err
}

func (c *Client) Get(url string) (*http.Response, error) {
	start := time.Now()
	resp, err := c.Client.Get(url)
	c.timit(start, resp, err)
	return resp, err
}

func (c *Client) Head(url string) (*http.Response, error) {
	start := time.Now()
	resp, err := c.Client.Head(url)
	c.timit(start, resp, err)
	return resp, err
}

func (c *Client) Post(url string, contentType string, body io.Reader) (*http.Response, error) {
	start := time.Now()
	resp, err := c.Client.Post(url, contentType, body)
	c.timit(start, resp, err)
	return resp, err
}

func (c *Client) PostForm(url string, data url.Values) (*http.Response, error) {
	start := time.Now()
	resp, err := c.Client.PostForm(url, data)
	c.timit(start, resp, err)
	return resp, err
}

func (c *Client) PostJson(url string, v interface{}) (*http.Response, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return c.Post(url, "application/json", bytes.NewReader(data))
}

func (c *Client) PostXml(url string, v interface{}) (*http.Response, error) {
	data, err := xml.Marshal(v)
	if err != nil {
		return nil, err
	}
	return c.Post(url, "application/xml", bytes.NewReader(data))
}
