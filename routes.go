package main

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/plimble/ace"
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
	db, err := sqlx.Connect("postgres", "user=dev dbname=webshop sslmode=disable")
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

	router.POST("/results", DisplaySearch)

	router.GET("/noresults", Noresults)

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
