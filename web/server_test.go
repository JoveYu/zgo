package web

import (
	"testing"

	"github.com/JoveYu/zgo/log"
)

func TestRouter(t *testing.T) {
	log.Install("stdout")

	server := NewServer()
	server.Router("GET", "^/$", thandler)
	server.Router("GET", "/(?P<name>\\w+)$", thandler)

	server.Run("127.0.0.1:7000")
}
