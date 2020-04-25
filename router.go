package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
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
		http.MethodHead,
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
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

// JankedHandler implements http.Handler
// F must be one of
// 	* func(ResponseWriter, *Request)
// 	* func(ResponseWriter, *Request, interface{})
// and must output one of
// 	* nil
// 	* interface{}
// 	* error
// 	* (interface{}, error)
type JankedHandler struct {
	F interface{}
}

// ServeHTTP - first figure out what F is and call it
func (hw JankedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// function itself
	fn := reflect.ValueOf(hw.F)
	// function signature
	sig := fn.Type()

	// check that the first two arguments are of the correct type
	if sig.NumIn() < 2 {
		log.Panicln("inner function has wrong amount of args")
	}
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
	// call the function
	rets := fn.Call(argv)
	// give the output to the ResponseWriter
	for _, r := range rets {
		if r.Kind() == reflect.Struct || r.Kind() == reflect.Slice {
			payload, _ := json.Marshal(r.Interface())
			w.Write(payload)
			break
		}
		if r.IsNil() {
			continue
		}
		if e := (*error)(nil); r.Type().Implements(reflect.TypeOf(e).Elem()) {
			msg := fmt.Sprintf("%v", r)
			payload, err := json.Marshal(struct{ Error string }{msg})
			if err != nil {
				fmt.Println(err.Error())
			}
			w.Write(payload)
		}
	}
}

// NewRouter ...
func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()
	for _, route := range Routes {
		handler := JankedHandler{route.HandlerFunc}
		withMiddleware := applyMiddleware(handler, route.Config, Middlewares...)
		mux.Handle(route.Pattern, withMiddleware)
	}
	return mux
}
