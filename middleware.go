package main

import (
	"log"
	"net/http"
	"time"
)

type middleware func(http.Handler, RouteConfig) http.Handler

func applyMiddleware(h http.Handler, cfg RouteConfig, mw ...middleware) http.Handler {
	for _, m := range mw {
		h = m(h, cfg)
	}
	return h
}

func requestRouteLogger(next http.Handler, c RouteConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf(
			"%s %s %s %s",
			r.Method,
			r.RequestURI,
			time.Since(start),
			getRequestID(r),
		)
	})
}

func requestIDGenerator(next http.Handler, cfg RouteConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nr := setRequestID(r)
		defer w.Header().Add("X-Request-ID", string(getRequestID(r)))
		next.ServeHTTP(w, nr)
	})
}

func requestAuthenticator(next http.Handler, cfg RouteConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authed := authChallenge(r)
		if authed {
			log.Printf("request %s authed", getRequestID(r))
		} else {
			log.Printf("request %s anonymous", getRequestID(r))
		}
		next.ServeHTTP(w, r)
	})
}
