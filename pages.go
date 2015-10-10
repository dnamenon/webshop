package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
)

// Page Handlers

func Home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	sess, err := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)

	err = sess.Delete("itemid")
	log.Println(err)
	ShowItemsHome(w, r)
}

func Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {


	sess, err := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)

	err = sess.Delete("itemid")
	log.Println(err)

	SimplePage(w, r, "try")

}

func Authfail(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	log.Print("Authorization Failed")
}

func UserPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	SimpleAuthenticatedPage(w,r, "user", nil)
}

func Itemid(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

	url := ps.ByName("id")
	str := strings.Trim(url, ":")
	fmt.Println(str)

	sess, err := globalSessions.SessionStart(w, req)
	defer sess.SessionRelease(w)

	sessname := sess.Get("itemid")
	log.Println(sessname)

	if err == nil {

		err = sess.Set("itemid", str)

		ShowItemid(w,req, str)
	}

}

func Signup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {


	sess, err := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)

	err = sess.Delete("itemid")
	log.Println(err)
	SimplePage(w, r, "signup")
}

func Pagination(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {


	sess, err := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)

	err = sess.Delete("itemid")

	url := ps.ByName("pg")
	str := strings.Trim(url, ":")

	var b int
	if _, err = fmt.Sscanf(str, "%5d", &b); err == nil {
		log.Println(err)
		ShowItemsPages(w,r, b)
	}
}

func DisplaySearch(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {


	sess, err := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)

	err = sess.Delete("itemid")

	url := ps.ByName("pg")
	str := strings.Trim(url, ":")

	var b int
	if _, err = fmt.Sscanf(str, "%5d", &b); err == nil {
		Search(w , r , b)
	}
}

func Noresults(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {


	sess, err := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)

	err = sess.Delete("itemid")
	log.Println(err)

	log.Println("login probz")
	SimplePage(w, r, "results")
}

func CartPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	Cart(w,r)
}

func BuyPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	Buy(w,r)
}

func BadPage(w http.ResponseWriter, r *http.Request, str string) {

	render := render.New(render.Options{})
	render.HTML(w, http.StatusOK, "badpage", str)
}
