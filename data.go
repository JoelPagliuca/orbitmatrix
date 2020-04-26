package main

import (
	"fmt"
	"log"
	"time"
)

// DB "Database"
type DB struct {
	items []Item
	users []User
}

// D database singleton
var D *DB

type model struct {
	ID          uuid
	DateCreated time.Time
}

func (m *model) init() {
	m.ID = generateUUID()
	m.DateCreated = time.Now()
}

// Item a configuration item
type Item struct {
	model
	Name        string
	Description string
}

// User ...
type User struct {
	model
	Name  string
	Token string `json:"-"`
}

// InitDB creates db object
func InitDB() {
	D = &DB{
		items: []Item{},
		users: []User{},
	}
}

// AddItem ...
func (db *DB) AddItem(i Item) (Item, error) {
	i.init()
	db.items = append(db.items, i)
	return i, nil
}

// GetItems ...
func (db *DB) GetItems() []Item {
	return db.items
}

// GetItem ...
func (db *DB) GetItem(id uint) (*Item, error) {
	if int(id) < len(db.items) {
		return &db.items[id], nil
	}
	return nil, fmt.Errorf("id not in database")
}

// AddUser ...
func (db *DB) AddUser(u User) (User, error) {
	u.init()
	u.Token = generateToken()
	db.users = append(db.users, u)
	log.Println("New user created: " + u.ID)
	return u, nil
}

// GetUserByToken ...
func (db *DB) GetUserByToken(t string) *User {
	for _, u := range db.users {
		if u.Token == t {
			return &u
		}
	}
	return nil
}
