package main

import (
	"appengine"
	"appengine/datastore"
	"appengine/user"

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
	e.Get("/api/v1.0/login", login)
	e.Get("/api/v1.0/logout", logout)

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
		return c.JSON(http.StatusInternalServerError, err)
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

func add(c *echo.Context) error {
	ac := appengine.NewContext(c.Request())

	m := member{}
	m.Events = []eventSignup{eventSignup{}, eventSignup{}}
	m.Office = []officeBearer{officeBearer{}, officeBearer{}}
	m.Contact.Phone = []phone{phone{}, phone{}}
	m.Contact.Address = []address{address{}, address{}}

	key, err := datastore.Put(ac, datastore.NewIncompleteKey(ac, "member", nil), &m)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	var m2 member
	if err = datastore.Get(ac, key, &m2); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, m2)
}
