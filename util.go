package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type key int

const (
	requestID   key = iota
	requestUser key = iota
)

func getEnv(key, defaultVal string) string {
	if value, found := os.LookupEnv(key); found {
		return value
	}
	log.Printf("Using default value for %s", key)
	return defaultVal
}

type uuid string

func generateUUID() uuid {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	u := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid(u)
}

func generateToken() string {
	crap := fmt.Sprintf("%s%s", generateUUID(), generateUUID())
	encoded := base64.StdEncoding.EncodeToString(bytes.NewBufferString(crap).Bytes())
	output := strings.TrimRight(encoded, "=")
	return output[0:50]
}

func getRequestID(req *http.Request) uuid {
	rid := req.Context().Value(requestID)
	if rid != nil {
		return rid.(uuid)
	}
	return ""
}

func setRequestID(req *http.Request) *http.Request {
	ctx := context.WithValue(req.Context(), requestID, generateUUID())
	return req.WithContext(ctx)
}

func getUser(req *http.Request) *User {
	u := req.Context().Value(requestUser)
	if u != nil {
		usr := u.(User)
		return &usr
	}
	return nil
}

// authChallenge check the api key on the request and attach the user to the context
// returns whether user is logged in
func authChallenge(req *http.Request) (*http.Request, bool) {
	authHeader := req.Header.Get("Authorization")
	if len(authHeader) == 0 {
		return nil, false
	}
	token := authHeader[7:]
	u := GetUserByToken(token)
	if u != nil {
		ctx := context.WithValue(req.Context(), requestUser, *u)
		return req.WithContext(ctx), true
	}
	return nil, false
}
