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

// GetTransitiveMemberOf get the groups the user is member of
func GetTransitiveMemberOf(db *gorm.DB, userID uint) []Group {
	grps := make(map[uint]Group)
	out := _getMemberOf(db, userID)
	for _, grp := range out {
		grps[grp.ID] = grp
		_getTransitiveMemberOf(db, grp.ID, grps)
	}
	out = []Group{}
	for _, g := range grps {

		out = append(out, g)
	}
	return out
}

// _getMemberOf immediate groups of the user
func _getMemberOf(db *gorm.DB, uid uint) []Group {
	grps := []Group{}
	gids := []uint{}
	rows, err := db.Table("group_users").Where("user_id = ?", uid).Select("group_id").Rows()
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		var gid uint
		err = rows.Scan(&gid)
		if err != nil {
			log.Println(err)
			return grps
		}
		gids = append(gids, gid)
	}
	for _, i := range gids {
		var grp Group
		if err := db.First(&grp, i).Error; err == nil {
			grps = append(grps, grp)
		}
	}
	return grps
}

func _getSubgroupOf(db *gorm.DB, gid uint) []Group {
	grps := []Group{}
	gids := []uint{}
	rows, err := db.Table("subgroups").Where("subgroup_id = ?", gid).Select("group_id").Rows()
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		var gid uint
		err = rows.Scan(&gid)
		if err != nil {
			log.Println(err)
			return grps
		}
		gids = append(gids, gid)
	}
	for _, gid := range gids {
		var grp Group
		if err := db.First(&grp, gid).Error; err == nil {
			grps = append(grps, grp)
		}
	}
	return grps
}

func _getTransitiveMemberOf(db *gorm.DB, gid uint, g map[uint]Group) {
	grps := _getSubgroupOf(db, gid)
	for _, grp := range grps {
		g[grp.ID] = grp
		_getTransitiveMemberOf(db, grp.ID, g)
	}
}
