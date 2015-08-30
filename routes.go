package main

import (
	"log"
	"time"

	"github.com/astaxie/beego/session"
	"github.com/jmoiron/sqlx"
	"github.com/plimble/ace"
)

var globalSessions *session.Manager

type Item struct {
	Id          int
	Title       string
	Date        time.Time
	Description string
	Seller      string
	Price       string
	Image       string
	Quantity    int
	Num         int
}

var db *sqlx.DB = SetupDB()

func SetupDB() *sqlx.DB {
	db, err := sqlx.Connect("postgres", "user=devmenon dbname=webshop sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	return db
}

func init() {
	globalSessions, _ = session.NewManager("file", `{"cookieName":"gosessionid","gclifetime":3600,"ProviderConfig":"./auth"}`)
	go globalSessions.GC()

}

func main() {

	defer db.Close()

	router := ace.Default()

	router.GET("/", Home)

	router.GET("/page/:pg", Pagination)

	router.POST("/results/:pg", DisplaySearch)

	router.GET("/noresults", Noresults)

	router.GET("/login", Login)

	router.POST("/login", PostLogin)

	router.GET("/signup", Signup)

	router.POST("/signup", PostSignup)

	router.GET("/logout", Logout)

	router.GET("/authfail", Authfail)

	router.GET("/user", UserPage)

	router.GET("/item/:id", Itemid)

	router.POST("/deletefromcart", DeleteFromCart)

	router.GET("/addtocart", AddToCart)

	router.GET("/cart", CartPage)

	router.GET("/buy", BuyPage)

	router.POST("/checkout", Checkout)

	router.GET("/public/css/:cssfile", fileLoadHandler)

	router.Run(":2500")
}
