package patch

import (
	"net/http"

	"github.com/JoveYu/zgo/http/httplog"
)

func init() {
	http.DefaultTransport = httplog.DefaultTransport
}
