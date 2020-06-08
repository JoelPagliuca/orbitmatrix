package main

import (
	"os"
	"strings"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func newTestDB(t *testing.T) *gorm.DB {
	os.Unsetenv("CALIBAN_DB_FILE")
	db := InitDB()
	t.Cleanup(func() {
		*db = *InitDB()
	})
	return db
}

func TestGetTransitiveMembers(t *testing.T) {
	db := newTestDB(t)
	g1 := Group{Name: "g1"}
	g2 := Group{Name: "g2"}
	g3 := Group{Name: "g3"}
	u1 := User{Name: "u1"}
	u2 := User{Name: "u2"}
	u3 := User{Name: "u3"}
	db.Create(&g1)
	db.Create(&g2)
	db.Create(&g3)
	db.Create(&u1)
	db.Create(&u2)
	db.Create(&u3)
	db.Model(&g1).Association("Subgroups").Append(&g2)
	db.Model(&g1).Association("Members").Append(&u1)
	db.Model(&g2).Association("Subgroups").Append(&g3)
	db.Model(&g2).Association("Members").Append(&u1)
	db.Model(&g2).Association("Members").Append(&u2)
	db.Model(&g3).Association("Members").Append(&u3)
	out := GetTransitiveMembers(db, g1.ID)
	outS := ""
	for _, u := range out {
		outS += u.Name
	}
	for _, u := range []User{u1, u2, u3} {
		if !strings.Contains(outS, u.Name) {
			t.Errorf("output missing " + u.Name)
		}
	}
	if len(out) != 3 {
		t.Error("output wasn't 3 users")
	}
}
