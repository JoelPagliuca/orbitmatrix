package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"os"
)

type key int

const (
	requestID key = iota
)

func getEnv(key, defaultVal string) string {
	if value, found := os.LookupEnv(key); found {
		return value
	}
	log.Printf("Using default value for %s", key)
	return defaultVal
}

func generateUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

func getRequestID(req *http.Request) string {
	rid := req.Context().Value(requestID)
	if rid != nil {
		return rid.(string)
	}
	return ""
}

func setRequestID(req *http.Request) *http.Request {
	ctx := context.WithValue(req.Context(), requestID, generateUUID())
	return req.WithContext(ctx)
}

func hydrateFromRequest(req *http.Request, thing interface{}) error { return nil }
