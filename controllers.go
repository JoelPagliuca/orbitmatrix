package main

import (
	"log"
	"net/http"
	"time"
)

func healthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func routeLogger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		log.Printf(
			"%s %s %s",
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	}
}
