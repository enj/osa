package osa

import (
	"time"

	"google.golang.org/appengine/datastore"
)

type Family struct {
	Key      *datastore.Key `json:"id" datastore:"-"`
	Name     string         `json:"name" endpoints:"req"`
	Keyword  string         `json:"keyword" endpoints:"req"`
	Email    string         `json:"email" endpoints:"req"`
	Login    string         `json:"login" endpoints:"req"`
	Password string         `json:"password" datastore:"-" endpoints:"req"`
	Hash     string         `json:"-"`
	Salt     string         `json:"-"`
	Active   bool           `json:"-"`
	Log      LogData        `json:"-"`
}

type Event struct {
	Key         *datastore.Key `json:"id" datastore:"-"`
	Title       string         `json:"title" endpoints:"req"`
	Description string         `json:"description" endpoints:"req"`
	Location    Location       `json:"location" endpoints:"req"`
	Duration    TimeRange      `json:"duration" endpoints:"req"`
	EnrollBy    time.Time      `json:"enrollBy" endpoints:"req"`
	Log         LogData        `json:"-"`
}

type Member struct {
	Key *datastore.Key `json:"id" datastore:"-"`
	//TODO figure out how to make ancestor
	//May be an issue since this can't be changed after creation
	//FamilyKey *datastore.Key `json:"familyID"`
	Primary      bool           `json:"-"`
	Active       bool           `json:"-"`
	Gender       string         `json:"gender,omitempty"`
	Relationship string         `json:"relationship" endpoints:"req"`
	Birthday     time.Time      `json:"birthday,omitempty"`
	Name         Name           `json:"name" endpoints:"req"`
	Contact      Contact        `json:"contact" endpoints:"req"`
	Office       []OfficeBearer `json:"-"`
	Log          LogData        `json:"-"`
}

type Name struct {
	Title  string `json:"title,omitempty"`
	First  string `json:"first" endpoints:"req"`
	Middle string `json:"middle,omitempty"`
	Last   string `json:"last" endpoints:"req"`
}

type Contact struct {
	Email struct {
		Primary     string `json:"primary" endpoints:"req"`
		Alternative string `json:"alternative,omitempty"`
	} `json:"email" endpoints:"req"`
	Phone   []Phone   `json:"phone,omitempty"`
	Address []Address `json:"address,omitempty"`
}

type Phone struct {
	Type struct {
		// work, home, etc
		Contact string `json:"contact" endpoints:"req"`
		// mobile, cell, fax, etc
		Phone string `json:"phone" endpoints:"req"`
	} `json:"type" endpoints:"req"`
	Number string `json:"number" endpoints:"req"`
}

type Address struct {
	Name     string   `json:"name" endpoints:"req"`
	Type     string   `json:"type" endpoints:"req"`
	Location Location `json:"location" endpoints:"req"`
}

type AddressLine struct {
	One   string `json:"one" endpoints:"req"`
	Two   string `json:"two,omitempty"`
	Three string `json:"three,omitempty"`
}

type Location struct {
	Line    AddressLine `json:"line" endpoints:"req"`
	City    string      `json:"city" endpoints:"req"`
	State   string      `json:"state" endpoints:"req"`
	Zip     string      `json:"zip" endpoints:"req"`
	Country string      `json:"country" endpoints:"req"`
}

type OfficeBearer struct {
	Type string    `json:"type" endpoints:"req"`
	Term TimeRange `json:"term" endpoints:"req"`
}

type TimeRange struct {
	Start time.Time `json:"start" endpoints:"req"`
	End   time.Time `json:"end" endpoints:"req"`
}

type EventSignup struct {
	EventKey *datastore.Key `json:"eventID" endpoints:"req"`
	Comments string         `json:"comments,omitempty" datastore:",noindex"`
	Time     time.Time      `json:"-"`
}

type LogData struct {
	Modified   time.Time      `json:"-"`
	ModifiedBy *datastore.Key `json:"-"`
}

type EventsList struct {
	Events []*Event `json:"events"`
}

type EventsListReq struct {
	Limit int `json:"limit" endpoints:"d=10,min=1,max=20,desc=The number of events to list."`
}

type EventsService struct {
}
