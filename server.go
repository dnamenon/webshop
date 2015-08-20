package main

import (
	"database/sql"
	"fmt"
	//"html/template"
	"log"
	"io/ioutil"
  "net/http"
  "path/filepath"
  //"path"
  //"reflect"
  "strings"
  "time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/plimble/ace"
	//"github.com/astaxie/beego/session"
	"github.com/unrolled/render"
	"golang.org/x/crypto/bcrypt"
	//"github.com/dnamenon/paginater"
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

func SetupDB() *sqlx.DB {
	db, err := sqlx.Connect("postgres", "user=devmenon dbname=webshop sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	return db
}



func main() {

	defer db.Close()

	router := ace.Default()
	router.Use(ace.Logger())

	router.GET("/", Home)

	router.GET("/page/:pg", Pagination)

	router.GET("/login", Login)

	router.POST("/login", PostLogin)

	router.GET("/signup", Signup)

	router.POST("/signup", PostSignup)

	router.GET("/logout", Logout)

	router.GET("/authfail", Authfail)

	router.GET("/user", User)

	router.GET("/item/:id", Itemid)

	router.GET("/public/css/:cssfile", fileLoadHandler)

	router.Run(":2500")
}

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

func Pagination(c *ace.C){
	var w = c.Writer
	var r = c.Request

	url := c.Param("pg")

	str := strings.Trim(url, ":")


	var b int
	if _, err := fmt.Sscanf(str, "%5d", &b); err == nil{
		ShowItemsPages(w, r, b)
	}
}
// end of page handlers
//action handlers begin

func SimplePage(w http.ResponseWriter, req *http.Request, template string) {

	r := render.New(render.Options{})
	r.HTML(w, http.StatusOK, template, nil)
}

func SimpleAuthenticatedPage(w http.ResponseWriter, req *http.Request, template string) {




	r := render.New(render.Options{})
	r.HTML(w, http.StatusOK, template, nil)
}

func PostLogin(c *ace.C) {
	var w = c.Writer
	var req = c.Request



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




	http.Redirect(w, req, "/user", 302)
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


	http.Redirect(w, req, "/", 302)
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
		render.HTML(w, http.StatusOK, "home", map[string]interface{} {
        "item0": items[0],
				"item1": items[1],
				"item2": items[2],
    })

}

func ShowItemsPages(w http.ResponseWriter, r *http.Request, num int) {

	// Loop through rows using only one struct
	items := []Item{}
	err := db.Select(&items, "SELECT * FROM items ORDER BY id ASC")
	if err != nil {
		  log.Println("Is this the problem?")
			log.Println(err)
			return
	}
factor := num-1

one,two,three := factor*3, (factor*3)+1, (factor*3)+2

log.Println("%d\n %d\n %d\n", one, two, three)

		render := render.New(render.Options{
			IndentJSON: true,
		})
		render.HTML(w, http.StatusOK, "home", map[string]interface{} {
        "item0": items[one],
				"item1": items[two],
				"item2": items[three],
    })

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

func fileLoadHandler(c *ace.C) {
	url := c.Param("cssfile")
	str := strings.TrimLeft(url, ":")
	w := c.Writer

 baseDir, _ := filepath.Abs("/Users/devmenon/golang/src/webshop/public/css/")

 fileName := "/"+str

 file, err := ioutil.ReadFile(fmt.Sprintf("%s%s", baseDir, fileName))

 w.Header().Set("Content-Type", "text/css")
 fmt.Fprint(w, string(file))
 if err != nil {
  fmt.Println("Error reading file: ")
  fmt.Println(err)
 }
}
