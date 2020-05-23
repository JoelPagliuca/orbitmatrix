package main

import (
	"log"
	"net/http"
)

var cfg struct {
	apiPort string
}

func init() {
	cfg.apiPort = getEnv("API_PORT", "7080")
	InitDB()
}

func main() {
	router := NewRouter()
	defer D.Close()
	log.Printf("Server started on " + cfg.apiPort)
	log.Fatal(http.ListenAndServe(":"+cfg.apiPort, router))
}
