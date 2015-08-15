package main

import (
	"database/sql"
	"net/http"
	"html/template"
	"log"
	"fmt"
	"path"
	"reflect"
  "time"

  "webshop/shopsess"

<<<<<<< HEAD

	"github.com/jmoiron/sqlx"
   "github.com/plimble/ace"
	_ "github.com/lib/pq"
=======
	"github.com/codegangsta/negroni"
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
>>>>>>> parent of 92090ac... back to httprouter and negroni-sessions
	"github.com/unrolled/render"
	"golang.org/x/crypto/bcrypt"
	"github.com/julienschmidt/httprouter"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Item struct {

	Id          int
	Title       string
	Date        time.Time
	Description string
	Seller      string
	Price       string
	Image       string
}





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

<<<<<<< HEAD

	router := ace.Default()

=======
	s := shopsess.Sessions()

  router := httprouter.New()
	n := negroni.Classic()
  n.UseHandler(router)
  n.Use(s)

	store := cookiestore.New([]byte("ohhhsooosecret"))
	n.Use(sessions.Sessions("global_session_store", store))
>>>>>>> parent of 92090ac... back to httprouter and negroni-sessions

  router.GET("/", http.HandlerFunc(Home))

<<<<<<< HEAD
	router.Get("/login",  http.HandlerFunc(Login))
=======
  router.GET("/login", Login)
>>>>>>> parent of 92090ac... back to httprouter and negroni-sessions

	router.HandleFunc("POST", "/login",  http.HandlerFunc(PostLogin))

<<<<<<< HEAD
	router.Get("/signup",  http.HandlerFunc(Signup))
=======
  router.GET("/signup", Signup)
>>>>>>> parent of 92090ac... back to httprouter and negroni-sessions

	router.HandleFunc("POST", "/signup",  http.HandlerFunc(PostSignup))

<<<<<<< HEAD
	router.HandleFunc("GET", "/logout",  http.HandlerFunc(Logout))
=======
	router.GET("/logout", LogoutPage)
>>>>>>> parent of 92090ac... back to httprouter and negroni-sessions

	router.Get("/authfail", http.HandlerFunc(Authfail))

	router.Get("/user", http.HandlerFunc(User))

	router.Get("/item/:id", http.HandlerFunc(Itemid))

  router.Run(":2500")
}

<<<<<<< HEAD
func Home(w http.ResponseWriter, r *http.Request) {

	ShowItems(w, r)

}

func Login(w http.ResponseWriter, r *http.Request) {
	SimplePage(w, r, "try")

}

func Authfail(w http.ResponseWriter, r *http.Request) {
	log.Print("Authorization Failed")
}

func User(w http.ResponseWriter, r *http.Request) {
	SimpleAuthenticatedPage(w, r, "user")
}

func Itemid(w http.ResponseWriter, r *http.Request) {
	url := hitch.Params(r).ByName("id")
	var str string = url
  fmt.Println(str)
	var b int
		if _, err := fmt.Sscanf(str, ":%s5d", &b); err == nil {
			ShowItemid(w, r, b)
		}

}

func Signup(w http.ResponseWriter, r *http.Request) {
	SimplePage(w, r, "signup")
=======

func Home(w http.ResponseWriter, r *http.Request,  _ httprouter.Params) {
ShowItems(w, r)

}

func Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	SimplePage(w,r, "try")

}

func LogoutPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		Logout(w, r)
	}

func Authfail(w http.ResponseWriter, r *http.Request,  _ httprouter.Params){
				log.Print("Authorization Failed")
}

func User(w http.ResponseWriter, r *http.Request,  _ httprouter.Params) {
SimplePage(w, r, "user")
}

func Itemid(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
   str := ps[0].Value

	 var b int
	 if _, err := fmt.Sscanf(str, ":%5d", &b); err == nil{
		 ShowItemid(w,r,b)
	 }

}

func Signup(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	SimplePage(w,r,"signup")
>>>>>>> parent of 92090ac... back to httprouter and negroni-sessions
}

// end of page handlers
//action handlers begin

func SimplePage(w http.ResponseWriter, req *http.Request, template string) {
	r := render.New(render.Options{})
	r.HTML(w, http.StatusOK, template, nil)
}


func SimpleAuthenticatedPage(w http.ResponseWriter, req *http.Request, template string) {
<<<<<<< HEAD
session := session.GetSession(req)

s := session.Id
=======
	session := sessions.GetSession(req)
	sess := session.Get("email")

	if sess == nil {
		http.Redirect(w, req, "/authfail", 301)
	}
>>>>>>> parent of 92090ac... back to httprouter and negroni-sessions

if s == nil{
	http.Redirect(w, req, "/authfail",301)
}
	r := render.New(render.Options{})
	r.HTML(w, http.StatusOK, template, nil)
}



func PostLogin(w http.ResponseWriter, req *http.Request){

var password_in_database string
var email string


username, password := req.FormValue("inputUsername"), req.FormValue("inputPassword")
err := db.QueryRow("SELECT email, password FROM users WHERE username=$1", username).Scan(&email, &password_in_database)
fmt.Println(password_in_database)
		if err == sql.ErrNoRows {
			http.Redirect(w, req, "/authfail", 301)
	} else if err != nil {
<<<<<<< HEAD
		log.Print(err)
		http.Redirect(w, req, "/authfail", 301)
	}



=======
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

fmt.Println(req)
>>>>>>> parent of 92090ac... back to httprouter and negroni-sessions

		session := sessions.GetSession(req)
		session.Set("email", email)
		http.Redirect(w, req, "/user", 302)
}



func PostSignup(w http.ResponseWriter, req *http.Request) {
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


func Logout(w http.ResponseWriter, req *http.Request) {

	http.Redirect(w, req, "/", 302)
}


func ShowItems(w http.ResponseWriter, r *http.Request) {
	// Loop through rows using only one struct
  item := Item{}
     rows, err := db.Queryx("SELECT * FROM items")
     for rows.Next() {
         err = rows.StructScan(&item)
         if err != nil {
					 log.Print("is this the problem?")
             log.Fatalln(err)
         }

  fmt.Printf("%#v\n", item)
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

		    if err := tmpl.Execute(w, item); err != nil {
					fmt.Println("is this the error")
		        http.Error(w, err.Error(), http.StatusInternalServerError)
		    }

     }

}

func ShowItemid(w http.ResponseWriter, r *http.Request, i int){
	itemid := Item{}
	    err := db.Get(&itemid, "SELECT * FROM items WHERE id=$1", i)
			if err != nil{
				log.Print("This must be the problem")
			}

fmt.Printf("%#v\n", itemid.Id)
if itemid.Id == 0{
	fmt.Println("My Probz")

	 http.Redirect(w, r, "/", 302)
}
render := render.New(render.Options{})
render.HTML(w, http.StatusOK, "item", &itemid)
}
