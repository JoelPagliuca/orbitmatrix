package main

import (
	"net/http"
)

// Route used to generate router
type Route struct {
	Pattern     string
	HandlerFunc interface{}
	Config      RouteConfig
}

// RouteConfig ...
type RouteConfig struct {
	AllowAnonymous bool
	AllowedMethods []string
}

var defaultRouteConfig = RouteConfig{
	AllowAnonymous: false,
	AllowedMethods: []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodHead,
	},
}

// Routes all the routes for the api
var Routes = []Route{
	Route{
		"/health",
		healthcheck,
		RouteConfig{
			AllowAnonymous: true,
			AllowedMethods: []string{"GET"},
		},
	},
	Route{
		"/user/me",
		getMe,
		defaultRouteConfig,
	},
	Route{
		"/user/register",
		registerUser,
		RouteConfig{
			AllowAnonymous: true,
			AllowedMethods: []string{"POST"},
		},
	},
	Route{
		"/item",
		getItems,
		defaultRouteConfig,
	},
	Route{
		"/item/add",
		addItem,
		defaultRouteConfig,
	},
}

// Middlewares the middleware to apply to all the above functions
var Middlewares = []middleware{
	requestMethodChecker,
	requestAuthenticator,
	requestRouteLogger,
	requestIDGenerator,
}

// NewRouter ...
func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()
	for _, route := range Routes {
		handler := JankedHandler{route.HandlerFunc}
		withMiddleware := applyMiddleware(handler, route.Config, Middlewares...)
		mux.Handle(route.Pattern, withMiddleware)
	}
	AddSwagger(mux, Routes)
	return mux
}
