package main

import (
	"database/sql"
	"net/http"
	"html/template"
	"log"
	"fmt"


	"github.com/codegangsta/negroni"
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	_ "github.com/lib/pq"
	"github.com/unrolled/render"
	"golang.org/x/crypto/bcrypt"

)

type Item struct {
	Id          int
	Title       string
	Date        string
	Description template.HTML
	Seller      string
	Price       float64
	Image       string
}

func (i *Item) id() int { return i.Id }
func (i *Item) title() string { return i.Title }
func (i *Item) date() string { return  i.Date }
func (i *Item) description() template.HTML { return i.Description }
func (i *Item) seller() string { return i.Seller}
func (i *Item) price() float64 { return i.Price}
func (i *Item) image() string { return i.Image}





var db *sql.DB = SetupDB()

func SetupDB() * sql.DB{
	db, err := sql.Open("postgres", "postgres://devmenon:cloud01@localhost:5432/webshop?sslmode=disable")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}

	return db
}



func main() {



	defer db.Close()


	mux := http.NewServeMux()
	n := negroni.Classic()
  n.UseHandler(mux)
	store := cookiestore.New([]byte("ohhhsooosecret"))
	n.Use(sessions.Sessions("global_session_store", store))


	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var items []Item

		rows, err := db.Query("select * from items")
		if err != nil {
			panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
		}
		defer rows.Close()

			for rows.Next() {
				var item Item
				var description string
				err := rows.Scan(&item.Id, &item.Title, &item.Date, &description, &item.Seller, &item.Price)
					if err != nil {
						panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
					}
		item.Description = template.HTML(description)
		items = append(items, item)

			}
			err = rows.Err()
			if err != nil {
				panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
			}

		ShowItems(w, r)
		SimplePage(w, r, "home")
	})


	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			SimplePage(w, r, "try")
		} else if r.Method == "POST" {
			fmt.Printf("Check")
			PostLogin(w, r)
		}
	})

	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
			Logout(w, r)
		})

		mux.HandleFunc("/authfail", func(w http.ResponseWriter, r *http.Request){
						log.Print("Authorization Failed")
		})

		mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		SimpleAuthenticatedPage(w, r, "user")
	})



	n.Run(":2500")
}



func SimplePage(w http.ResponseWriter, req *http.Request, template string) {
	r := render.New(render.Options{})
	r.HTML(w, http.StatusOK, template, nil)
}


func SimpleAuthenticatedPage(w http.ResponseWriter, req *http.Request, template string) {
	session := sessions.GetSession(req)
	sess := session.Get("email")

	if sess == nil {
		http.Redirect(w, req, "/authfail", 301)
	}

	r := render.New(render.Options{})
	r.HTML(w, http.StatusOK, template, nil)
}



func errHandler(err error) {
	if err != nil {
		log.Print(err)
	}
}


func PostLogin(w http.ResponseWriter, req *http.Request){


session := sessions.GetSession(req)

var password_in_database string
var email string


	username, password := req.FormValue("inputUsername"), req.FormValue("inputPassword")
	err := db.QueryRow("SELECT email, password FROM users WHERE username=$1", username).Scan(&email, &password_in_database)

		if err == sql.ErrNoRows {
			http.Redirect(w, req, "/authfail", 301)
	} else if err != nil {
			log.Print(err)
			http.Redirect(w, req, "/authfail", 301)
		}

	err = bcrypt.CompareHashAndPassword([]byte(password_in_database), []byte(password))
		if err == bcrypt.ErrMismatchedHashAndPassword {
			http.Redirect(w, req, "/authfail", 301)
		} else if err != nil {
			log.Print(err)
			http.Redirect(w, req, "/authfail", 301)
		}

		session.Set("email", email)
		http.Redirect(w, req, "/user", 302)
}



func PostSignup(w http.ResponseWriter, req *http.Request) {
	username = req.FormValue("inputUsername")
	password = req.FormValue("inputPassword")
	email = req.FormValue("inputEmail")

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


func Logout(w http.ResponseWriter, req *http.Request) {
	session := sessions.GetSession(req)
	session.Delete("email")



		http.Redirect(w, req, "/", 302)
}



func ShowItems(w http.ResponseWriter, r *http.Request) {
    item := Item{}

    fp := path.Join("templates", "home.tmpl")
    tmpl, err := template.ParseFiles(fp)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if err := tmpl.Execute(w, item); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
