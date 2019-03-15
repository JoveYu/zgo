// context for web framework

package web

import (
	"context"
	"encoding/json"
	"fmt"
	"mime"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"strings"
)

type ContextHandlerFunc func(Context)

type LoggingResponseWriter struct {
	http.ResponseWriter
	status int
}

type Context struct {
	context.Context
	Request        *http.Request
	ResponseWriter *LoggingResponseWriter
	Charset        string

	// for router params
	Params map[string]string

	formParsed bool
}

func NewContext(w http.ResponseWriter, r *http.Request) Context {
	return Context{
		Request:        r,
		ResponseWriter: &LoggingResponseWriter{w, 0},
		Charset:        "utf-8",
		formParsed:     false,
	}
}

func ContextHandler(f ContextHandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := NewContext(w, r)
		f(ctx)
	})
}

func ContextCancelHandler(f ContextHandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := NewContext(w, r)
		c, cancel := context.WithCancel(context.Background())
		defer cancel()

		ctx.Context = c
		f(ctx)
	})
}

func (w *LoggingResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (ctx *Context) Method() string {
	return ctx.Request.Method
}

func (ctx *Context) URL() *url.URL {
	return ctx.Request.URL
}

func (ctx *Context) WriteHeader(status int) {
	ctx.ResponseWriter.WriteHeader(status)
}

func (ctx *Context) WriteString(s string) {
	ctx.ResponseWriter.Write([]byte(s))
}

func (ctx *Context) WriteJSON(v interface{}) error {
	ctx.SetContextType("application/json")
	return json.NewEncoder(ctx.ResponseWriter).Encode(v)
}

func (ctx *Context) WriteFile(path string) {
	http.ServeFile(ctx.ResponseWriter, ctx.Request, path)
}

func (ctx *Context) Headers() http.Header {
	return ctx.Request.Header
}

func (ctx *Context) GetHeader(k string) string {
	return ctx.Request.Header.Get(k)
}

func (ctx *Context) SetHeader(k string, v string) {
	ctx.ResponseWriter.Header().Set(k, v)
}

func (ctx *Context) AddHeader(k string, v string) {
	ctx.ResponseWriter.Header().Add(k, v)
}

func (ctx *Context) Cookies() []*http.Cookie {
	return ctx.Request.Cookies()
}

func (ctx *Context) GetCookie(k string) *http.Cookie {
	c, err := ctx.Request.Cookie(k)
	if err != nil {
		return nil
	}
	return c
}

func (ctx *Context) GetCookieV(k string) string {
	c := ctx.GetCookie(k)
	if c != nil {
		return c.Value
	}
	return ""
}

func (ctx *Context) SetCookie(c *http.Cookie) {
	http.SetCookie(ctx.ResponseWriter, c)
}

func (ctx *Context) SetCookieKV(k string, v string) {
	c := http.Cookie{
		Name:     k,
		Value:    v,
		Path:     "/",
		HttpOnly: true,
	}
	ctx.SetCookie(&c)
}

func (ctx *Context) DelCookie(k string) {
	c := http.Cookie{
		Name:   k,
		Path:   "/",
		MaxAge: -1,
	}
	ctx.SetCookie(&c)
}

func (ctx *Context) UserAgent() string {
	return ctx.GetHeader("User-Agent")
}

func (ctx *Context) Query() url.Values {
	return ctx.Request.URL.Query()
}

func (ctx *Context) GetQuery(k string) string {
	return ctx.Request.URL.Query().Get(k)
}

func (ctx *Context) FormFile(k string) (multipart.File, *multipart.FileHeader, error) {
	return ctx.Request.FormFile(k)
}

func (ctx *Context) Form() url.Values {
	if !ctx.formParsed {
		ctx.Request.ParseForm()
		ctx.formParsed = true
	}
	return ctx.Request.Form
}

func (ctx *Context) GetForm(k string) string {
	if !ctx.formParsed {
		ctx.Request.ParseForm()
		ctx.formParsed = true
	}
	return ctx.Request.Form.Get(k)
}

func (ctx *Context) SetCharset(c string) {
	ctx.Charset = c
}

// allow use ext to set content type
// SetContextType("json")
// SetContextType("application/json")
func (ctx *Context) SetContextType(t string) string {
	// if is ext
	if !strings.ContainsRune(t, '/') {
		t = mime.TypeByExtension(fmt.Sprintf(".%s", t))
	}
	if t != "" {
		ctx.SetHeader("Context-Type", fmt.Sprintf("%s; charset=%s", t, ctx.Charset))
	}
	return t
}

func (ctx *Context) ClientIP() string {
	clientIP := ctx.GetHeader("X-Forwarded-For")
	clientIP = strings.TrimSpace(strings.Split(clientIP, ",")[0])
	if clientIP == "" {
		clientIP = strings.TrimSpace(ctx.GetHeader("X-Real-Ip"))
	}
	if clientIP != "" {
		return clientIP
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(ctx.Request.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}

func (ctx *Context) Abort(status int, body string) {
	ctx.SetContextType("text/plain")
	ctx.WriteHeader(status)
	ctx.WriteString(body)
}

func (ctx *Context) AbortJSON(status int, v interface{}) {
	ctx.WriteHeader(status)
	ctx.WriteJSON(v)
}

func (ctx *Context) Redirect(status int, url string) {
	ctx.SetHeader("Location", url)
	ctx.WriteHeader(status)
	ctx.WriteString("Redirecting to ")
	ctx.WriteString(url)
}
