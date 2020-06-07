package main

import (
	"os"

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
	Members     []User `gorm:"many2many:group_users;"`
}

// InitDB creates db object
func InitDB() {
	filename, ok := os.LookupEnv("CALIBAN_DB_FILE")
	if !ok {
		filename = ":memory:"
	}
	DB, err := gorm.Open("sqlite3", filename)
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

// GetUserByToken ...
func GetUserByToken(t string) *User {
	tok := Token{}
	err := D.Preload("User").Where("value = ?", t).First(&tok).Error
	if err != nil {
		return nil
	}
	return &tok.User
}

// TODO: GetTransitiveMembers
func GetTransitiveMembers(groupID uint) []User {
	return []User{}
}

// TODO: GetTransitiveMemberOf
func GetTransitiveMemberOf(userID uint) []Group {
	return []Group{}
}
