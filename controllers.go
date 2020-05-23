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
	D.Select("1")
	if len(D.GetErrors()) == 0 {
		out.Entries.Database = "Healthy"
	} else {
		log.Printf("DB Errors: %s", D.Error)
	}
	return out
}

// ITEM

func getItems(w http.ResponseWriter, r *http.Request) []Item {
	items := []Item{}
	D.Find(&items)
	return items
}

type addItemInput struct {
	I Item `from:"body"`
}

func addItem(w http.ResponseWriter, r *http.Request, in addItemInput) Item {
	D.Create(&in.I)
	return in.I
}

// USER

func getMe(w http.ResponseWriter, r *http.Request) User {
	u := getUser(r)
	return *u
}

type tokenResponse struct {
	ID        uint
	TokenType string
	Token     string
}

type registerUserInput struct {
	U User `from:"body"`
}

func registerUser(w http.ResponseWriter, r *http.Request, in registerUserInput) tokenResponse {
	if !D.NewRecord(in.U) {
		return tokenResponse{}
	}
	D.Create(&in.U)
	t := Token{UserID: in.U.ID}
	D.Create(&t)
	res := tokenResponse{
		ID:        in.U.ID,
		TokenType: "Bearer",
		Token:     t.Value,
	}
	return res
}

// GROUP

type createGroupInput struct {
}

func createGroup(w http.ResponseWriter, r *http.Request, in createGroupInput) {

}
