package web

import (
	"regexp"
)

type Router struct {
	r       string
	cr      *regexp.Regexp
	method  string
	handler ContextHandlerFunc
}
