package main

import (
	"encoding/json"
	"log"
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

func addItem(w http.ResponseWriter, r *http.Request) {
	var newItem Item
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&newItem)
	if err != nil {
		log.Println("Error: ", err.Error())
	}
	out, _ := D.AddItem(newItem)
	payload, _ := json.Marshal(out)
	w.Write(payload)
}
