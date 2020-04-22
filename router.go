package main

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
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

// Middlewares the middleware to apply to all the above functions
var Middlewares = []middleware{
	requestRouteLogger,
	requestIDGenerator,
}

// HandlerWrapper implements http.Handler
type HandlerWrapper struct {
	// must be a func(ResponseWriter, *Request)
	F interface{}
}

func (hw HandlerWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// function itself
	fn := reflect.ValueOf(hw.F)
	// function signature
	sig := fn.Type()
	if sig.NumIn() < 2 {
		log.Panicln("inner function has wrong amount of args")
	}
	// check that the first two arguments are of the correct type
	if !(reflect.TypeOf(w).AssignableTo(sig.In(0)) && reflect.TypeOf(r).AssignableTo(sig.In(1))) {
		log.Panicln("First two args of a function were bad")
	}
	argv := make([]reflect.Value, sig.NumIn())
	argv[0] = reflect.ValueOf(w)
	argv[1] = reflect.ValueOf(r)
	// be helpful for the rest of the inut args
	for i := 2; i < sig.NumIn(); i++ {
		// fill up a struct arg from the request body
		if sig.In(i).Kind() == reflect.Struct {
			arg := reflect.New(sig.In(i)).Interface()
			defer r.Body.Close()
			err := json.NewDecoder(r.Body).Decode(&arg)
			if err != nil {
				log.Println("Error: ", err.Error())
			}
			argv[i] = reflect.ValueOf(arg).Elem()
		}
	}
	fn.Call(argv)
}

// NewRouter ...
func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()
	for _, route := range Routes {
		handler := HandlerWrapper{route.HandlerFunc}
		withMiddleware := applyMiddleware(handler, Middlewares...)
		mux.Handle(route.Pattern, withMiddleware)
	}
	return mux
}
