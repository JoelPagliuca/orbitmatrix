package main

import (
	"log"
	"net/http"
)

func healthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

// ITEM

func getItems(w http.ResponseWriter, r *http.Request) []Item {
	items := D.GetItems()
	return items
}

func addItem(w http.ResponseWriter, r *http.Request, newItem Item) Item {
	out, _ := D.AddItem(newItem)
	return out
}

// USER

func getMe(w http.ResponseWriter, r *http.Request) User {
	u := getUser(r)
	return *u
}

type tokenResponse struct {
	ID        uuid
	TokenType string
	Token     string
}

func registerUser(w http.ResponseWriter, r *http.Request, u User) tokenResponse {
	log.Println(u.Name)
	u, _ = D.AddUser(u)
	res := tokenResponse{
		ID:        u.ID,
		TokenType: "Bearer",
		Token:     u.Token,
	}
	return res
}
