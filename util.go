package main

import (
	"log"
	"os"
)

func getEnv(key, defaultVal string) string {
	if value, found := os.LookupEnv(key); found {
		return value
	}
	log.Printf("Using default value for %s", key)
	return defaultVal
}
