package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/plimble/ace"
	"github.com/unrolled/render"
)

// Page Handlers

func Home(c *ace.C) {
	var w = c.Writer
	var r = c.Request

	sess, err := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)

	err = sess.Delete("itemid")
	log.Println(err)
	ShowItemsHome(w, r)
}

func Login(c *ace.C) {
	var w = c.Writer
	var r = c.Request

	sess, err := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)

	err = sess.Delete("itemid")
	log.Println(err)

	SimplePage(w, r, "try")

}

func Authfail(c *ace.C) {

	log.Print("Authorization Failed")
}

func UserPage(c *ace.C) {

	SimpleAuthenticatedPage(c, "user", nil)
}

func Itemid(c *ace.C) {
	req := c.Request
	w := c.Writer

	url := c.Param("id")
	str := strings.Trim(url, ":")
	fmt.Println(str)

	sess, err := globalSessions.SessionStart(w, req)
	defer sess.SessionRelease(w)

	sessname := sess.Get("itemid")
	log.Println(sessname)

	if err == nil {

		err = sess.Set("itemid", str)

		ShowItemid(c, str)
	}

}

func Signup(c *ace.C) {
	var w = c.Writer
	var r = c.Request

	sess, err := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)

	err = sess.Delete("itemid")
	log.Println(err)
	SimplePage(w, r, "signup")
}

func Pagination(c *ace.C) {
	var w = c.Writer
	var r = c.Request

	sess, err := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)

	err = sess.Delete("itemid")

	url := c.Param("pg")
	str := strings.Trim(url, ":")

	var b int
	if _, err = fmt.Sscanf(str, "%5d", &b); err == nil {
		log.Println(err)
		ShowItemsPages(c, b)
	}
}

func DisplaySearch(c *ace.C) {
	var w = c.Writer
	var r = c.Request

	sess, err := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)

	err = sess.Delete("itemid")

	url := c.Param("pg")
	str := strings.Trim(url, ":")

	var b int
	if _, err = fmt.Sscanf(str, "%5d", &b); err == nil {
		Search(c, b)
	}
}

func Noresults(c *ace.C) {
	var w = c.Writer
	var r = c.Request

	sess, err := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)

	err = sess.Delete("itemid")
	log.Println(err)

	log.Println("login probz")
	SimplePage(w, r, "results")
}

func CartPage(c *ace.C) {

	Cart(c)
}

func BuyPage(c *ace.C) {

	Buy(c)
}

func BadPage(c *ace.C, str string) {
	var w = c.Writer

	render := render.New(render.Options{})
	render.HTML(w, http.StatusOK, "badpage", str)
}
