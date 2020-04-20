package main

import (
	"log"
	"net/http"
	"time"
)

func requestRouteLogger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		log.Printf(
			"%s %s %s %s",
			r.Method,
			r.RequestURI,
			time.Since(start),
			getRequestID(r),
		)
	}
}

func requestIDGenerator(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nr := setRequestID(r)
		next(w, nr)
	}
}
