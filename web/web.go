package web

import ()

var (
	DefaultServer = NewServer()
)

func Route(method string, path string, f ContextHandlerFunc) {
	DefaultServer.Router(method, path, f)
}

func GET(path string, f ContextHandlerFunc) {
	DefaultServer.Router("GET", path, f)
}

func POST(path string, f ContextHandlerFunc) {
	DefaultServer.Router("POST", path, f)
}

func PUT(path string, f ContextHandlerFunc) {
	DefaultServer.Router("PUT", path, f)
}

func DELETE(path string, f ContextHandlerFunc) {
	DefaultServer.Router("DELETE", path, f)
}

func PATCH(path string, f ContextHandlerFunc) {
	DefaultServer.Router("PATCH", path, f)
}

func Run(addr string) {
	DefaultServer.Run(addr)
}
