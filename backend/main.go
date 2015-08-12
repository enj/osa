package main

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"

	"github.com/labstack/echo"

	"net/http"
)

var (
	e = wrapMux()
)

// some middleware / handlers
func wrapMux() *echo.Echo {

	// Echo instance
	e := echo.New()

	// Add CORS middleware

	// Routes
	e.Get("/api/v1.0/add", add)
	e.Get("/api/v1.0/members", members)
	e.Get("/api/v1.0/events", allEvents)
	e.Get("/api/v1.0/user_events", userEvents)
	e.Get("/api/v1.0/login", login)
	e.Get("/api/v1.0/logout", logout)
	e.Post("/api/v1.0/event_signup", memberEventSignup)

	return e
}

// Handler
func members(c *echo.Context) error {

	ac := appengine.NewContext(c.Request())

	// The Query type and its methods are used to construct a query.
	q := datastore.NewQuery("member")

	// To retrieve the results,
	// you must execute the Query using its GetAll or Run methods.
	var people []member
	_, err := q.GetAll(ac, &people)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responseJSON{"", err.Error()})
	}

	return c.JSON(http.StatusOK, people)
}

func login(c *echo.Context) error {
	ac := appengine.NewContext(c.Request())
	u := user.Current(ac)
	redirect := c.Query("continue")
	if redirect == "" {
		redirect = "/"
	}
	if u == nil {
		url, _ := user.LoginURL(ac, redirect)
		return c.Redirect(http.StatusFound, url)
	}
	return c.Redirect(http.StatusFound, redirect)
}

func logout(c *echo.Context) error {
	ac := appengine.NewContext(c.Request())
	u := user.Current(ac)

	if u == nil {
		return c.Redirect(http.StatusFound, "/")
	}

	url, _ := user.LogoutURL(ac, "/")
	return c.Redirect(http.StatusFound, url)
}

func init() {
	http.Handle("/api/", e)
}

func memberEventSignup(c *echo.Context) error {
	// TODO redo this to run in a transaction since it needs to be atomic
	eventName := c.Form("event")
	if eventName == "" {
		return c.JSON(http.StatusBadRequest, responseJSON{"", "Missing event name"})
	}
	ac := appengine.NewContext(c.Request())
	e := event{}
	eventKey := datastore.NewKey(ac, "event", eventName, 0, nil)
	if err := datastore.Get(ac, eventKey, &e); err != nil {
		return c.JSON(http.StatusBadRequest, responseJSON{"", err.Error()})
	}
	// TODO do more validation here like seeing if event is still active
	comments := c.Form("comments")
	m, k := getOrCreateMember(ac)
	q := datastore.NewQuery("member").
		Filter("__key__ =", k).
		Filter("Events.Event =", eventName)
	count, err := q.Count(ac)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responseJSON{"", err.Error()})
	}
	alreadySignedUp := count != 0 //TODO redo data model to fix this?
	if alreadySignedUp {
		return c.JSON(http.StatusBadRequest, responseJSON{"", "Already signed up"})
	} else {
		eventDetails := eventSignup{e.Title, comments, time.Now()}
		m.Events = append(m.Events, eventDetails)
		_, err := datastore.Put(ac, k, &m)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, responseJSON{"", err.Error()})
		}
		return c.JSON(http.StatusOK, responseJSON{"Signed up for event", ""})
	}
}

func allEvents(c *echo.Context) error {
	ac := appengine.NewContext(c.Request())

	// The Query type and its methods are used to construct a query.
	q := datastore.NewQuery("event")

	// To retrieve the results,
	// you must execute the Query using its GetAll or Run methods.
	var events []event
	_, err := q.GetAll(ac, &events)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responseJSON{"", err.Error()})
	}

	return c.JSON(http.StatusOK, events)
}

func userEvents(c *echo.Context) error {
	ac := appengine.NewContext(c.Request())
	m, _ := getOrCreateMember(ac)
	if m.Events == nil || len(m.Events) == 0 {
		c.JSON(http.StatusOK, responseJSON{"", "no events"})
	}
	return c.JSON(http.StatusOK, m.Events)
}

func getOrCreateMember(ac context.Context) (member, *datastore.Key) {
	var m member
	u := user.Current(ac)
	id := u.ID
	key := datastore.NewKey(ac, "member", id, 0, nil)
	err := datastore.Get(ac, key, &m)
	if err != nil {
		m2 := member{}
		m2.Account = id
		m2.Contact.Email.Primary = u.Email
		m2.Modified = time.Now()
		datastore.Put(ac, key, &m2) //TODO check for error?
		return m2, key
	}
	return m, key
}

func add(c *echo.Context) error {
	ac := appengine.NewContext(c.Request())

	//m := member{}
	//m.Events = []eventSignup{eventSignup{}, eventSignup{}}
	//m.Office = []officeBearer{officeBearer{}, officeBearer{}}
	//m.Contact.Phone = []phone{phone{}, phone{}}
	//m.Contact.Address = []address{address{}, address{}}

	//key, err := datastore.Put(ac, datastore.NewIncompleteKey(ac, "member", nil), &m)
	//if err != nil {
	//	return c.JSON(http.StatusInternalServerError, responseJSON{"", err.Error()})
	//}

	//var m2 member
	//if err = datastore.Get(ac, key, &m2); err != nil {
	//	return c.JSON(http.StatusInternalServerError, responseJSON{"", err.Error()})
	//}

	//return c.JSON(http.StatusOK, m2)

	e := event{}
	e.Title = "Awesome Event"
	key := datastore.NewKey(ac, "event", e.Title, 0, nil)
	_, err := datastore.Put(ac, key, &e)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responseJSON{"", err.Error()})
	}
	return c.JSON(http.StatusOK, e)
}
