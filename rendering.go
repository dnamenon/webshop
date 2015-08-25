package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/plimble/ace"
	"github.com/unrolled/render"
)

func SimplePage(w http.ResponseWriter, req *http.Request, template string) {

	r := render.New(render.Options{})
	r.HTML(w, http.StatusOK, template, nil)
}

func SimpleAuthenticatedPage(c *ace.C, template string, data interface{}) {
	w := c.Writer
	r := c.Request

	sess, err := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(c.Writer)

	sessname := sess.Get("email")
	log.Println(sessname)

	if err != nil {
		log.Println(err)
	} else {
		r := render.New(render.Options{})
		r.HTML(c.Writer, http.StatusOK, template, data)
	}
}

func fileLoadHandler(c *ace.C) {
	url := c.Param("cssfile")
	str := strings.TrimLeft(url, ":")
	w := c.Writer

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
