package router

import (
	"net/http"

	"github.com/serbi/nasa_apod_collector/internal/pkg/models"
)

type RouteWrapper func(r models.Route) models.Route

type localRoute struct {
	path    string
	handler models.HandlerFunc
}

func (l localRoute) Handler() models.HandlerFunc {
	return l.handler
}

func (l localRoute) Path() string {
	return l.path
}

func NewRoute(path string, handler models.HandlerFunc, opts ...RouteWrapper) models.Route {
	var r models.Route = localRoute{path, handler}
	for _, o := range opts {
		r = o(r)
	}
	return r
}

type localRouter struct {
	routes []models.Route
}

func (lr *localRouter) Routes() []models.Route {
	return lr.routes
}

func (lr *localRouter) initRoutes() {
	lr.routes = []models.Route{
		NewRoute(PicturesPath, HandlePictures),
	}
}

func (lr *localRouter) RegisterHandlers() {
	for _, route := range lr.Routes() {
		http.HandleFunc(route.Path(), route.Handler())
	}
}

func NewRouter() (lr *localRouter) {
	lr = &localRouter{}
	lr.initRoutes()
	return
}
