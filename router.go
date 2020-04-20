package main

import (
	"net/http"
)

// Route used to generate router
type Route struct {
	Pattern     string
	HandlerFunc http.HandlerFunc
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

// NewRouter ...
func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()
	for _, route := range Routes {
		handler := route.HandlerFunc
		handler = requestRouteLogger(handler)
		handler = requestIDGenerator(handler)
		mux.Handle(route.Pattern, handler)
	}
	return mux
}
