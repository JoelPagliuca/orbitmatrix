package main

import (
	"log"
	"net/http"
	"strings"
	"time"
)

type middleware func(http.Handler, RouteConfig) http.Handler

func applyMiddleware(h http.Handler, cfg RouteConfig, mw ...middleware) http.Handler {
	for _, m := range mw {
		h = m(h, cfg)
	}
	return h
}

func requestIDGenerator(next http.Handler, cfg RouteConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nr := setRequestID(r)
		w.Header().Add("X-Request-ID", string(getRequestID(nr)))
		next.ServeHTTP(w, nr)
	})
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

func requestAuthenticator(next http.Handler, cfg RouteConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, authed := authChallenge(r)
		if !authed {
			log.Printf("request %s anonymous", getRequestID(r))
			if cfg.AllowAnonymous {
				next.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
			}
			return
		}
		next.ServeHTTP(w, req)
	})
}

func requestMethodChecker(next http.Handler, cfg RouteConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.Header().Add("Allow", strings.Join(cfg.AllowedMethods, ", "))
			w.WriteHeader(http.StatusNoContent)
			return
		}
		for _, m := range cfg.AllowedMethods {
			if r.Method == m {
				next.ServeHTTP(w, r)
				return
			}
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})
}
