package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/plimble/ace"
	//"github.com/astaxie/beego/session"
	"github.com/unrolled/render"

	"golang.org/x/crypto/bcrypt"
)

//action handlers begin

func PostLogin(c *ace.C) {
	var w = c.Writer
	var req = c.Request

	sess, _ := globalSessions.SessionStart(w, req)
	defer sess.SessionRelease(w)
	sessname := sess.Get("email")
	cartsess := sess.Get("cart")
	var cart []string
	fmt.Println(sessname)
	fmt.Println(cartsess)

	var password_in_database string
	var email string

	badstr := "Authorization Failed"

	username, password := req.FormValue("inputUsername"), req.FormValue("inputPassword")
	err := db.QueryRow("SELECT email, password FROM users WHERE username=$1", username).Scan(&email, &password_in_database)
	if err == sql.ErrNoRows {
		BadPage(c, badstr)
	} else if err != nil {
		log.Print(err)
		BadPage(c, badstr)
	}

	err = bcrypt.CompareHashAndPassword([]byte(password_in_database), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		BadPage(c, badstr)
	} else if err != nil {
		log.Print(err)
		BadPage(c, badstr)
	}

	if err == nil {

		err2 := sess.Delete("email")
		err2 = sess.Set("email", email)
		err2 = sess.Set("cart", cart)
		log.Println(req)
		if err2 == nil {
			c.Redirect("/")
		}
	}

}

func PostSignup(c *ace.C) {
	var w = c.Writer
	var req = c.Request

	username := req.FormValue("reg_username")
	password := req.FormValue("reg_password")
	email := req.FormValue("reg_email")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("INSERT INTO users (username, password, email) VALUES ($1, $2, $3)", username, string(hashedPassword), email)
	if err != nil {
		log.Print(err)
	}

	http.Redirect(w, req, "/login", 302)
}

func Logout(c *ace.C) {
	var w = c.Writer
	var req = c.Request

	sess, _ := globalSessions.SessionStart(w, req)
	globalSessions.SessionDestroy(w, req)
	err := sess.Delete("email")

	if err == nil {
		http.Redirect(w, req, "/", 302)
	}
}

func ShowItemsHome(w http.ResponseWriter, r *http.Request) {

	// Loop through rows using only one struct
	items := []Item{}
	err := db.Select(&items, "SELECT * FROM items ORDER BY id ASC")
	if err != nil {
		log.Println("Is this the problem?")
		log.Println(err)
		return
	}

	log.Println("%#v\n", items)

	render := render.New(render.Options{
		IndentJSON: true,
	})
	render.HTML(w, http.StatusOK, "home", map[string]interface{}{
		"item0": items[0],
		"item1": items[1],
		"item2": items[2],
	})

}

func ShowItemsPages(c *ace.C, i int) {
	var w = c.Writer

	// Loop through rows using only one struct
	items := []Item{}
	err := db.Select(&items, "SELECT * FROM items ORDER BY id ASC")
	if err != nil {
		log.Println("Is this the problem?")
		log.Println(err)
		return
	}

	max := len(items)
	maxparam := max - 3
	factor := i - 1
	one, two, three := factor*3, (factor*3)+1, (factor*3)+2

	log.Println("%d\n %d\n %d\n", one, two, three)
	log.Println(max)
	render := render.New(render.Options{
		IndentJSON: true,
	})
	if one >= maxparam {
		if max%3 == 0 {
			render.HTML(w, http.StatusOK, "home", map[string]interface{}{
				"item0": items[one],
				"item1": items[two],
				"item2": items[three],
			})
		} else if max%3 == 2 {
			render.HTML(w, http.StatusOK, "home", map[string]interface{}{
				"item0": items[one],
				"item1": items[two],
			})
		} else if max%3 == 1 {
			render.HTML(w, http.StatusOK, "home", map[string]interface{}{
				"item0": items[one],
			})
		}

	} else if one < maxparam {
		render.HTML(w, http.StatusOK, "home", map[string]interface{}{
			"item0": items[one],
			"item1": items[two],
			"item2": items[three],
		})
	} else {
		c.Redirect("/")
	}

}

func ShowItemid(c *ace.C, s string) {
	var w = c.Writer
	var r = c.Request
	fmt.Println(s)
	itemid := Item{}
	err := db.Get(&itemid, "SELECT * FROM items WHERE id=$1", s)
	if err != nil {
		log.Print("This must be the problem")
	}

	if itemid.Id == 0 {

		http.Redirect(w, r, "/", 302)
	}

	render := render.New(render.Options{
		IndentJSON: true,
	})
	render.HTML(c.Writer, http.StatusOK, "item", &itemid)
}
