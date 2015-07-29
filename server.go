package main

import (
	"database/sql"
	"net/http"
	"html/template"
	"log"
	"fmt"
	"path"
	"reflect"




	"github.com/codegangsta/negroni"
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/julienschmidt/httprouter"
	"github.com/jmoiron/sqlx"
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





var db *sqlx.DB = SetupDB()

func SetupDB() *sqlx.DB{
	db, err := sqlx.Connect("postgres", "user=devmenon dbname=webshop sslmode=disable")
     if err != nil {
         log.Fatalln(err)
     }


	return db
}



func main() {



	defer db.Close()


//	mux := http.NewServeMux()
router := httprouter.New()
	n := negroni.Classic()
  n.UseHandler(router)

	store := cookiestore.New([]byte("ohhhsooosecret"))
	n.Use(sessions.Sessions("global_session_store", store))


	router.GET("/", Home)


	router.GET("/login", Login)

	router.GET("/logout", LogoutPage)

		router.GET("/authfail", Authfail)

		router.GET("/user", User)


	n.Run(":2500")
}


func Home(w http.ResponseWriter, r *http.Request,  _ httprouter.Params) {
//ShowItems(w, r)
SimplePage(w,r, "home")
}

func Login(w http.ResponseWriter, r *http.Request,  _ httprouter.Params) {
	if r.Method == "GET" {
		SimplePage(w, r, "try")
	} else if r.Method == "POST" {
		fmt.Println("Check")
		PostLogin(w, r)
	}
}

func LogoutPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		Logout(w, r)
	}

func Authfail(w http.ResponseWriter, r *http.Request,  _ httprouter.Params){
				log.Print("Authorization Failed")
}

func User(w http.ResponseWriter, r *http.Request,  _ httprouter.Params) {
SimpleAuthenticatedPage(w, r, "user")
}

// end of page handlers
//action handlers begin

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
	username := req.FormValue("inputUsername")
	password := req.FormValue("inputPassword")
	email := req.FormValue("inputEmail")

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
	// Loop through rows using only one struct
     item := Item{}
     rows, err := db.Queryx("SELECT * FROM items")
     for rows.Next() {
         err := rows.StructScan(&item)
         if err != nil {
					 log.Print("is this the problem?")
             log.Fatalln(err)
         }
         fmt.Printf("%#v\n", item)
     }


 fp := path.Join("templates", "home.tmpl")

 loop := reflect.ValueOf(item)

	 values := make([]interface{}, loop.NumField())

	 for i := 0; i < loop.NumField(); i++ {
			 values[i] = loop.Field(i).Interface()
	 }


		    tmpl, err := template.ParseFiles(fp)
		    if err != nil {

		        http.Error(w, err.Error(), http.StatusInternalServerError)
		        return
		    }

		    if err := tmpl.Execute(w, loop); err != nil {
		        http.Error(w, err.Error(), http.StatusInternalServerError)
		    }

}