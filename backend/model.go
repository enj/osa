package main

import (
	"time"
)

type family struct {
	Name    string
	Keyword string
	Email   string
	Login   string
	Hash    string
	Active  bool
	Log     logData
}

type event struct {
	Title       string
	Description string
	Location    location
	Duration    timeRange
	Enroll      time.Time
	Log         logData
}

type member struct {
	Account      string
	Family       string
	Primary      bool
	Active       bool
	Gender       bool
	Relationship string
	Birthday     time.Time
	Modified     time.Time
	Name         name
	Contact      contact
	Events       []eventSignup
	Office       []officeBearer
}

type name struct {
	Title  string
	First  string
	Middle string
	Last   string
}

type contact struct {
	Email struct {
		Primary     string
		Alternative string
	}
	Phone   []phone
	Address []address
}

type phone struct {
	Type struct {
		// work, home, etc
		Contact string
		// mobile, cell, fax, etc
		Phone string
	}
	Number string
}

type address struct {
	Name     string
	Type     string
	Location location
}

type location struct {
	Line struct {
		One string
		Two string
	}
	City    string
	State   string
	Zip     string
	Country string
}

type officeBearer struct {
	Type string
	Term timeRange
}

type timeRange struct {
	Start time.Time
	End   time.Time
}

type eventSignup struct {
	Event    string
	Comments string
	Time     time.Time
}

type logData struct {
	Modified   time.Time
	ModifiedBy string
}

type responseJSON struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}
