package web

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"regexp"
	"strings"
	"time"

	"github.com/JoveYu/zgo/log"
)

const (
	PoweredBy string = "zgo/0.0.1"
)

type Server struct {
	Addr    string
	Routers []Router
	Charset string
	Debug   bool
}

func NewServer() *Server {
	server := Server{
		Charset: "utf-8",
		Debug:   false,
	}
	return &server
}

func (s *Server) StaticFile(path string, dir string) {
	r := fmt.Sprintf("^%s.*$", path)
	handler := func(ctx Context) {
		// disable list directory
		if strings.HasSuffix(ctx.URL().Path, "/") {
			http.NotFound(ctx.ResponseWriter, ctx.Request)
			return
		}

		// XXX status 200 in log is wrong
		handler := http.StripPrefix(path, http.FileServer(http.Dir(dir)))
		handler.ServeHTTP(ctx.ResponseWriter, ctx.Request)
	}
	s.Router("GET", r, handler)
}

func (s *Server) Router(method string, path string, handlers ...ContextHandlerFunc) {
	cr, err := regexp.Compile(path)
	if err != nil {
		log.Warn("can not add route [%s] %s", path, err)
		return
	}

	s.Routers = append(s.Routers, Router{
		r:        path,
		cr:       cr,
		method:   method,
		handlers: handlers,
	})
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tstart := time.Now()

	ctx := NewContext(w, r)
	defer s.LogRequest(tstart, &ctx)

	// debug
	if s.Debug {
		ctx.Debug = true
		ctx.DebugBody = &strings.Builder{}

		// debug req
		data, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Error("can not dump req: %s", err)
		}
		for _, b := range strings.Split(string(data), "\n") {
			log.Debug("> %s", b)
		}

		// debug resp
		defer func(ctx Context) {
			log.Debug("< %s %d %s", ctx.Request.Proto,
				ctx.Flag.Status, http.StatusText(ctx.Flag.Status),
			)
			for k, v := range ctx.ResponseWriter.Header() {
				for _, vv := range v {
					log.Debug("< %s: %s", k, vv)
				}
				// XXX Content-Length and Date is missing
			}
			log.Debug("<")
			for _, b := range strings.Split(ctx.DebugBody.String(), "\n") {
				log.Debug("< %s", b)
			}
		}(ctx)
	}

	path := ctx.URL().Path

	// default header
	ctx.SetHeader("X-Powered-By", PoweredBy)
	ctx.SetContentType("text/plain")

	for _, router := range s.Routers {

		// HEAD request use GET Handler
		if ctx.Method() != router.method && !(ctx.Method() == "HEAD" && router.method == "GET") {
			continue
		}

		if !router.cr.MatchString(path) {
			continue
		}

		match := router.cr.FindStringSubmatch(path)
		if len(match[0]) != len(path) {
			continue
		}

		if len(match) > 1 {
			for idx, name := range router.cr.SubexpNames()[1:] {
				ctx.Params[name] = match[idx+1]
			}
		}

		for _, h := range router.handlers {
			h(ctx)
			if ctx.Flag.BreakNext {
				break
			}
		}
		return
	}

	ctx.Abort(http.StatusNotFound, http.StatusText(http.StatusNotFound))
}

func (s *Server) Run(addr string) error {
	return http.ListenAndServe(addr, s)
}

func (s *Server) LogRequest(tstart time.Time, ctx *Context) {

	log.Info("%d|%s|%s|%s|%s|%d",
		ctx.Flag.Status, ctx.Method(), ctx.URL().Path,
		ctx.Query().Encode(), ctx.ClientIP(),
		time.Since(tstart)/time.Microsecond,
	)
}
