package main

import (
	"log"
	"net/http"

	"github.com/drone/go.stripe"
	"github.com/plimble/ace"
	"github.com/unrolled/render"
)

func CartHandler(c *ace.C) {
	w := c.Writer
	r := c.Request

	sess, _ := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(c.Writer)

	sessname := sess.Get("email")
	log.Println(sessname)

	sessnamestr := sessname.(string)

	if sessname == nil {
		c.Redirect("/")
	} else {

		user := User{}
		err := db.Get(&user, "SELECT * FROM users WHERE email=$1", sessnamestr)
		if err != nil {
			log.Print("This must be the problem")
		}

		r := render.New(render.Options{})
		r.HTML(c.Writer, http.StatusOK, "cart", user)
	}
}

func Buy(c *ace.C) {
	stripe.SetKey("sk_test_s7zyOcXwo4E9YLlBXpL4Ie44")

}
