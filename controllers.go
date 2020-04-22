package main

import (
	"net/http"
)

func healthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func getItems(w http.ResponseWriter, r *http.Request) []Item {
	items := D.GetItems()
	return items
}

func addItem(w http.ResponseWriter, r *http.Request, newItem Item) Item {
	out, _ := D.AddItem(newItem)
	return out
}
