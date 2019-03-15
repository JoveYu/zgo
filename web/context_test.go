package web

import (
	"net/http"
	"testing"

	"github.com/JoveYu/zgo/log"
)

func thandler(ctx Context) {
	m := ctx.Method()
	url := ctx.URL()
	log.Debug("%s %s", m, url)
	a := ctx.GetQuery("a")
	log.Debug("query a=%s", a)
	b := ctx.GetForm("b")
	log.Debug("form b=%s", b)
	ua := ctx.UserAgent()
	log.Debug("ua %s", ua)

	ctx.WriteHeader(200)
	ctx.WriteString("hello world")
}

func TestCtxHandler(t *testing.T) {
	log.Install("stdout")

	http.Handle("/", ContextHandler(thandler))
	http.ListenAndServe("127.0.0.1:7000", nil)
}
