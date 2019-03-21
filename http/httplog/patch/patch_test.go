package patch

import (
	"github.com/JoveYu/zgo/log"
	"net/http"
	"testing"
)

func TestPatch(t *testing.T) {
	log.Install("stdout")
	http.Get("http://httpbin.org/get")
}
