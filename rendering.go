package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
)

func SimplePage(w http.ResponseWriter, req *http.Request, template string) {

	r := render.New(render.Options{})
	r.HTML(w, http.StatusOK, template, nil)
}

func SimpleAuthenticatedPage(w http.ResponseWriter, r *http.Request, template string, data interface{}) {

	sess, err := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)

	sessname := sess.Get("email")
	log.Println(sessname)

	if err != nil {
		log.Println(err)
	} else {
		r := render.New(render.Options{})
		r.HTML(w, http.StatusOK, template, data)
	}
}

func fileLoadHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	url := ps.ByName("cssfile")
	str := strings.TrimLeft(url, ":")


	baseDir, _ := filepath.Abs("/Users/devmenon/golang/src/webshop/public/css/")

	fileName := "/" + str

	file, err := ioutil.ReadFile(fmt.Sprintf("%s%s", baseDir, fileName))

	w.Header().Set("Content-Type", "text/css")
	fmt.Fprint(w, string(file))
	if err != nil {
		fmt.Println("Error reading file: ")
		fmt.Println(err)
	}
}
