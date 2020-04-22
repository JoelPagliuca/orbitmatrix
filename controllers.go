package main

import (
	"encoding/json"
	"net/http"
)

func healthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func getItems(w http.ResponseWriter, r *http.Request) {
	items := D.GetItems()
	payload, _ := json.Marshal(items)
	w.Write(payload)
}

func addItem(w http.ResponseWriter, r *http.Request, newItem Item) {
	out, _ := D.AddItem(newItem)
	payload, _ := json.Marshal(out)
	w.Write(payload)
}
