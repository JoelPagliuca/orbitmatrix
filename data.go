package main

import (
	"fmt"
	"time"
)

// DB "Database"
type DB struct {
	data []Item
}

// D database singleton
var D *DB

// Item a configuration item
type Item struct {
	ID          uuid
	DateCreated time.Time
	Name        string
	Description string
}

// User ...
type User struct {
	ID          uuid
	DateCreated time.Time
	Name        string
	Key         string `json:"-"`
}

// Check make sure the model is allowed
func (i Item) Check() bool {
	return true
}

// InitDB creates db object
func InitDB() {
	D = &DB{
		data: []Item{},
	}
}

// newID gets the next available ID for a type
func (db *DB) newID() uuid {
	return generateUUID()
}

// AddItem ...
func (db *DB) AddItem(i Item) (Item, error) {
	i.ID = db.newID()
	i.DateCreated = time.Now()
	db.data = append(db.data, i)
	return i, nil
}

// GetItems ...
func (db *DB) GetItems() []Item {
	return db.data
}

// GetItem ...
func (db *DB) GetItem(id uint) (Item, error) {
	if int(id) < len(db.data) {
		return db.data[id], nil
	}
	return Item{}, fmt.Errorf("id not in database")
}
