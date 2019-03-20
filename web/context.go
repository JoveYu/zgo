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

type Context struct {
	context.Context
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	Charset        string

	// for debug
	Debug     bool
	DebugBody *strings.Builder

	// for router params
	Params map[string]string

	formParsed bool
	breakNext  bool
	status     int
}

func NewContext(w http.ResponseWriter, r *http.Request) Context {
	return Context{
		Request:        r,
		ResponseWriter: w,
		Charset:        "utf-8",

		// for debug
		Debug: false,

		Params: map[string]string{},

		formParsed: false,
		breakNext:  false,
		status:     200,
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

func (ctx *Context) BreakNext() {
	ctx.breakNext = true
}

func (ctx *Context) Param(k string) string {
	v, ok := ctx.Params[k]
	if !ok {
		return ""
	}
	return v
}

func (ctx *Context) Method() string {
	return ctx.Request.Method
}

func (ctx *Context) URL() *url.URL {
	return ctx.Request.URL
}

func (ctx *Context) ReadJSON(v interface{}) error {
	return json.NewDecoder(ctx.Request.Body).Decode(v)
}

func (ctx *Context) Write(b []byte) (int, error) {
	if ctx.Debug {
		ctx.DebugBody.Write(b)
	}

	return ctx.ResponseWriter.Write(b)
}

func (ctx *Context) WriteHeader(status int) {
	ctx.status = status
	ctx.ResponseWriter.WriteHeader(status)
}

func (ctx *Context) WriteString(s string) {
	ctx.Write([]byte(s))
}

func (ctx *Context) WriteJSON(v interface{}) error {
	ctx.SetContentType("application/json")
	ctx.WriteHeader(200)
	return json.NewEncoder(ctx).Encode(v)
}

func (ctx *Context) WriteJSONP(v interface{}) error {
	callback := ctx.GetQuery("callback")
	if callback != "" {
		ctx.SetContentType("application/javascript")
		ctx.WriteHeader(200)
		ctx.WriteString(fmt.Sprintf("%s(", callback))

		// XXX if err, body is wrong
		err := json.NewEncoder(ctx).Encode(v)
		if err != nil {
			return err
		}

		ctx.WriteString(")")
		return nil
	} else {
		return ctx.WriteJSON(v)
	}
}

func (ctx *Context) WriteFile(path string) {
	http.ServeFile(ctx.ResponseWriter, ctx.Request, path)
}

// simple cors allow ajax
func (ctx *Context) CORS() {
	origin := ctx.GetHeader("Origin")
	if origin != "" {
		ctx.SetHeader("Access-Control-Allow-Origin", origin)
		ctx.SetHeader("Access-Control-Allow-Credentials", "true")
	}

	method := ctx.GetHeader("Access-Control-Request-Method")
	if method != "" {
		ctx.SetHeader("Access-Control-Allow-Methods", method)
	}

	header := ctx.GetHeader("Access-Control-Request-Headers")
	if header != "" {
		ctx.SetHeader("Access-Control-Allow-Headers", header)
	}
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
// SetContentType("json")
// SetContentType("application/json")
func (ctx *Context) SetContentType(t string) string {
	// if is ext
	if !strings.ContainsRune(t, '/') {
		t = mime.TypeByExtension(fmt.Sprintf(".%s", t))
	}
	if t != "" {
		ctx.SetHeader("Content-Type", fmt.Sprintf("%s; charset=%s", t, ctx.Charset))
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
	ctx.SetContentType("text/plain")
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
