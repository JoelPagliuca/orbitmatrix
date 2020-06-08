package main

import (
	"log"
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
	Members     []*User  `gorm:"many2many:group_users;"`
	Subgroups   []*Group `gorm:"many2many:subgroups;association_jointable_foreignkey:subgroup_id"`
}

// InitDB creates db object
func InitDB() *gorm.DB {
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
	if DB.Error != nil {
		log.Println("Migration failed: " + DB.Error.Error())
	}
	D = DB
	return D
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

// GetTransitiveMembers gets members in this group including subgroups
func GetTransitiveMembers(db *gorm.DB, groupID uint) []User {
	users := make(map[uint]User)
	_getTransitiveMembers(db, groupID, users)
	out := []User{}
	for _, u := range users {
		out = append(out, u)
	}
	return out
}

func _getTransitiveMembers(db *gorm.DB, g uint, u map[uint]User) {
	group := Group{}
	if err := db.Preload("Subgroups").Preload("Members").First(&group, g).Error; err != nil {
		return
	}
	for _, usr := range group.Members {
		u[usr.ID] = *usr
	}
	for _, gr := range group.Subgroups {
		_getTransitiveMembers(db, gr.ID, u)
	}
}

// TODO: GetTransitiveMemberOf
func GetTransitiveMemberOf(db *gorm.DB, userID uint) []Group {
	return []Group{}
}
