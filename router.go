package main

import (
	"fmt"
	"net/http"
)

// Route used to generate router
type Route struct {
	Pattern     string
	HandlerFunc interface{}
	Config      *RouteConfig
}

// RouteConfig ...
type RouteConfig struct{}

// Routes all the routes for the api
var Routes = []Route{
	Route{
		"/health",
		healthcheck,
		nil,
	},
	Route{
		"/items",
		getItems,
		nil,
	},
	Route{
		"/items/add",
		addItem,
		nil,
	},
}

type HandlerWrapper struct {
	F interface{}
}

func (hw HandlerWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("wrapper instead of real func")
}

// NewRouter ...
func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()
	for _, route := range Routes {
		handler := HandlerWrapper{route.HandlerFunc}
		withMiddleware := requestRouteLogger(handler)
		withMiddleware = requestIDGenerator(handler)
		mux.Handle(route.Pattern, withMiddleware)
	}
	return mux
}
