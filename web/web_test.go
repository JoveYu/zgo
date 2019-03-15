package web

import (
	"testing"

	"github.com/JoveYu/zgo/log"
)

func TestWeb(t *testing.T) {
	log.Install("stdout")

	GET("^/$", thandler)
	GET("/(?P<name>\\w+)$", thandler)

	Run("127.0.0.1:7000")
}
