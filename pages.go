package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/plimble/ace"
)

// Page Handlers

func Home(c *ace.C) {
	var w = c.Writer
	var r = c.Request

	ShowItemsHome(w, r)
}

func Login(c *ace.C) {
	var w = c.Writer
	var r = c.Request

	SimplePage(w, r, "try")

}

func Authfail(c *ace.C) {

	log.Print("Authorization Failed")
}

func User(c *ace.C) {
	var w = c.Writer
	var r = c.Request

	SimpleAuthenticatedPage(w, r, "user")
}

func Itemid(c *ace.C) {

	url := c.Param("id")
	str := strings.Trim(url, ":")
	fmt.Println(str)

	ShowItemid(c, str)

}

func Signup(c *ace.C) {
	var w = c.Writer
	var r = c.Request

	SimplePage(w, r, "signup")
}

func Pagination(c *ace.C) {

	url := c.Param("pg")
	str := strings.Trim(url, ":")

	var b int
	if _, err := fmt.Sscanf(str, "%5d", &b); err == nil {
		ShowItemsPages(c, b)
	}
}

func DisplaySearch(c *ace.C) {
	Search(c)
}

func Noresults(c *ace.C) {
	var w = c.Writer
	var r = c.Request
	SimplePage(w, r, "results")
}
