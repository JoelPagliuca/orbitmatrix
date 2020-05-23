package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// D the database
var D *gorm.DB

// Item a configuration item
type Item struct {
	gorm.Model
	Name        string
	Description string
}

// User ...
type User struct {
	gorm.Model
	Name string
}

// Token ...
type Token struct {
	gorm.Model
	UserID uint
	User   User
	Value  string
}

// BeforeSave set the token value
func (t *Token) BeforeSave() (err error) {
	t.Value = generateToken()
	return nil
}

// Group ...
type Group struct {
	gorm.Model
	Name        string
	Description string
}

// InitDB creates db object
func InitDB() {
	DB, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		panic("Couldn't connect to database")
	}
	DB.AutoMigrate(
		&User{},
		&Token{},
		&Group{},
		&Item{},
	)
	D = DB
}

// GetUserByToken TODO
func GetUserByToken(t string) *User {
	tok := Token{}
	err := D.Preload("User").Where("value = ?", t).First(&tok).Error
	if err != nil {
		return nil
	}
	return &tok.User
}
