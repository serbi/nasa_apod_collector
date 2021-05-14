package models

import "net/http"

type HandlerFunc func(res http.ResponseWriter, req *http.Request)

type Router interface {
	Routes() []Route
}

type Route interface {
	Path() string
	Handler() HandlerFunc
}
