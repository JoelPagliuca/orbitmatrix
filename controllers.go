package main

import (
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

type addItemInput struct {
	I Item `from:"body"`
}

func addItem(w http.ResponseWriter, r *http.Request, in addItemInput) Item {
	out, _ := D.AddItem(in.I)
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

type registerUserInput struct {
	U User `from:"body"`
}

func registerUser(w http.ResponseWriter, r *http.Request, in registerUserInput) tokenResponse {
	u, _ := D.AddUser(in.U)
	res := tokenResponse{
		ID:        u.ID,
		TokenType: "Bearer",
		Token:     u.Token,
	}
	return res
}
