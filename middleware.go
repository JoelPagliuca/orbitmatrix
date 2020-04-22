package main

import (
	"log"
	"net/http"
	"time"
)

type middleware func(http.Handler) http.Handler

func applyMiddleware(h http.Handler, mw ...middleware) http.Handler {
	for _, m := range mw {
		h = m(h)
	}
	return h
}

func requestRouteLogger(next http.Handler) http.Handler {
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

func requestIDGenerator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nr := setRequestID(r)
		next.ServeHTTP(w, nr)
	})
}
