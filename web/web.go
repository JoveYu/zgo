package web

import (
	"net/http"
)

var (
	DefaultServer = NewServer()
)

func Route(method string, path string, f ...ContextHandlerFunc) {
	DefaultServer.Router(method, path, f...)
}

func GET(path string, f ...ContextHandlerFunc) {
	DefaultServer.Router(http.MethodGet, path, f...)
}

func POST(path string, f ...ContextHandlerFunc) {
	DefaultServer.Router(http.MethodPost, path, f...)
}

func PUT(path string, f ...ContextHandlerFunc) {
	DefaultServer.Router(http.MethodPut, path, f...)
}

func DELETE(path string, f ...ContextHandlerFunc) {
	DefaultServer.Router(http.MethodDelete, path, f...)
}

func PATCH(path string, f ...ContextHandlerFunc) {
	DefaultServer.Router(http.MethodPatch, path, f...)
}

func OPTIONS(path string, f ...ContextHandlerFunc) {
	DefaultServer.Router(http.MethodOptions, path, f...)
}

func StaticFile(path string, dir string) {
	DefaultServer.StaticFile(path, dir)
}

func Run(addr string) error {
	return DefaultServer.Run(addr)
}
