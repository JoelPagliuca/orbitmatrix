package main

import (
	"log"
	"net/http"
)

type healthResponse struct {
	Status  string
	Entries struct {
		Database string
	}
}

func healthcheck(w http.ResponseWriter, r *http.Request) healthResponse {
	out := healthResponse{
		Status: "Healthy",
	}
	D.Select("1")
	if len(D.GetErrors()) == 0 {
		out.Entries.Database = "Healthy"
	} else {
		log.Printf("DB Errors: %s", D.Error)
	}
	return out
}

// ITEM

func getItems(w http.ResponseWriter, r *http.Request) []Item {
	items := []Item{}
	D.Find(&items)
	return items
}

type itemInput struct {
	I Item `from:"body"`
}

func addItem(w http.ResponseWriter, r *http.Request, in itemInput) Item {
	if err := D.Create(&in.I).Error; err != nil {
		log.Println(err)
		return Item{}
	}
	return in.I
}

// USER

func getMe(w http.ResponseWriter, r *http.Request) User {
	u := getUser(r)
	return *u
}

type tokenResponse struct {
	ID        uint
	TokenType string
	Token     string
}

type registerUserInput struct {
	U User `from:"body"`
}

func registerUser(w http.ResponseWriter, r *http.Request, in registerUserInput) tokenResponse {
	if !D.NewRecord(in.U) {
		return tokenResponse{}
	}
	D.Create(&in.U)
	t := Token{UserID: in.U.ID}
	D.Create(&t)
	res := tokenResponse{
		ID:        in.U.ID,
		TokenType: "Bearer",
		Token:     t.Value,
	}
	return res
}

// GROUP

type groupInput struct {
	G Group `from:"body"`
}

func createGroup(w http.ResponseWriter, r *http.Request, in groupInput) Group {
	if err := D.Create(&in.G).Error; err != nil {
		log.Println(err)
		return Group{}
	}
	return in.G
}

type getGroupsInput struct {
	GroupID uint `from:"query"`
}

func getGroups(w http.ResponseWriter, r *http.Request, in getGroupsInput) []Group {
	groups := []Group{}
	if in.GroupID > 0 {
		D.Preload("Members").Find(&groups, in.GroupID)
	} else {
		D.Preload("Members").Find(&groups)
	}
	return groups
}

type addGroupMembersInput struct {
	GroupID    uint   `from:"query"`
	UserID     []uint `from:"query"`
	SubgroupID []uint `from:"query"`
}

func addGroupMembers(w http.ResponseWriter, r *http.Request, in addGroupMembersInput) {
	group := Group{}
	D.First(&group, in.GroupID)
	for _, u := range in.UserID {
		us := User{}
		us.ID = u
		D.Model(&group).Association("Members").Append(us)
	}
	for _, g := range in.SubgroupID {
		gr := Group{}
		gr.ID = g
		D.Model(&group).Association("Subgroups").Append(gr)
	}
	err := D.Error
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type getTransitiveMembersInput struct {
	GroupID uint `from:"query"`
}

func getTransitiveMembers(w http.ResponseWriter, r *http.Request, in getTransitiveMembersInput) []User {
	usrs := GetTransitiveMembers(D, in.GroupID)
	return usrs
}

type getTransitiveGroupsInput struct {
	UserID uint `from:"query"`
}

func getTransitiveGroups(w http.ResponseWriter, r *http.Request, in getTransitiveGroupsInput) []Group {
	grps := GetTransitiveMemberOf(D, in.UserID)
	return grps
}
