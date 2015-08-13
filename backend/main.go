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
	var members []member
	_, err := q.GetAll(ac, &members)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responseJSON{"", err.Error()})
	}

	// This will return a JSON null if there are no members
	return c.JSON(http.StatusOK, members)
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
	eventKey := datastore.NewKey(ac, "event", eventName, 0, nil)

	// TODO do more validation here like seeing if event is still active
	q := datastore.NewQuery("event").
		Filter("__key__ =", eventKey)

	count, err := q.Count(ac)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responseJSON{"", err.Error()})
	}
	if count != 1 {
		return c.JSON(http.StatusBadRequest, responseJSON{"", "Invalid event"})
	}

	comments := c.Form("comments")

	// TODO determine best time to cause creation of "member"
	//m := member{}
	//k, err := getOrCreateMember(ac, &m)
	//if err != nil {
	//	return c.JSON(http.StatusInternalServerError, responseJSON{"", err.Error()})
	//}

	// Use this member's key as parent of the event signup (even if the member does not yet exist)
	k, _ := getOrCreateMember(ac, nil) // Error checking not needed when m is nil

	eventSignupKey := datastore.NewKey(ac, "event_signup", eventName, 0, k)
	q = datastore.NewQuery("event_signup").
		Ancestor(k).
		Filter("__key__ =", eventSignupKey)
	count, err = q.Count(ac)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responseJSON{"", err.Error()})
	}

	if count == 1 {
		return c.JSON(http.StatusBadRequest, responseJSON{"", "Already signed up"})
	} else {
		eventSignupDetails := eventSignup{comments, time.Now()}
		_, err = datastore.Put(ac, eventSignupKey, &eventSignupDetails)
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
	k, _ := getOrCreateMember(ac, nil) //Don't need to check for DB errors when m is nil

	q := datastore.NewQuery("event_signup").Ancestor(k).KeysOnly()
	eventSignupKeys, err := q.GetAll(ac, nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responseJSON{"", err.Error()})
	}

	var eventKeys []*datastore.Key
	for _, eventSignupKey := range eventSignupKeys {
		eventKeys = append(eventKeys, datastore.NewKey(ac, "event", eventSignupKey.StringID(), 0, nil))
	}

	events := make([]event, len(eventKeys))
	if err = datastore.GetMulti(ac, eventKeys, events); err != nil {
		return c.JSON(http.StatusInternalServerError, responseJSON{"", err.Error()})
	}

	if events == nil || len(events) == 0 {
		return c.JSON(http.StatusOK, responseJSON{"no events", ""})
	}

	return c.JSON(http.StatusOK, events)
}

func getOrCreateMember(ac context.Context, m *member) (key *datastore.Key, err error) {
	u := user.Current(ac)
	key = datastore.NewKey(ac, "member", u.ID, 0, nil)

	// If no member is provided, then we only cared about having the key
	if m == nil {
		return
	}

	if err = datastore.Get(ac, key, m); err != nil {
		// Assumes the m is not modified when key does not exist
		m.Contact.Email.Primary = u.Email
		m.Modified = time.Now()
		_, err = datastore.Put(ac, key, m)
	}

	return
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
