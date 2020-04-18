package main

import (
	"net/http"
)

// Route used to generate router
type Route struct {
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
	Config      *RouteConfig
}

// RouteConfig ...
type RouteConfig struct{}

// Routes all the routes for the api
var Routes = []Route{
	Route{
		"GET",
		"/health",
		healthcheck,
		nil,
	},
}

// NewRouter ...
func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()
	for _, route := range Routes {
		handler := route.HandlerFunc
		handler = routeLogger(handler)
		mux.Handle(route.Pattern, handler)
	}
	return mux
}
