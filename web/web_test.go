package web

import (
	"testing"

	"github.com/JoveYu/zgo/log"
)

func handler1(ctx Context) {
	log.Debug("handler1")
	ctx.BreakNext()
}
func handler2(ctx Context) {
	log.Debug("handler2")
}

func TestWeb(t *testing.T) {
	log.Install("stdout")

	GET("^/$", handler1, handler2)
	GET("/(?P<name>\\w+)$", thandler)

	DefaultServer.Debug = true
	Run("127.0.0.1:7000")
}
