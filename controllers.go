package main

import (
	"log"
	"net/http"
)

type healthResponse struct {
	Status  string
	Entries struct {
		Database string
	}
}

func healthcheck(w http.ResponseWriter, r *http.Request) healthResponse {
	out := healthResponse{
		Status: "Healthy",
	}
	out.Entries.Database = "Healthy"
	return out
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
