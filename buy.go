package main

import (
	"log"
	"net/http"

	"github.com/drone/go.stripe"
	"github.com/plimble/ace"
	"github.com/unrolled/render"
)

func AddToCart(c *ace.C) {
	var w = c.Writer
	var r = c.Request

	sess, err := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)

	user := sess.Get("email")
	item := sess.Get("itemid")
	cart := sess.Get("cart")
	if cartstr, ok := cart.([]string); ok {
		if itemstr, ok2 := item.(string); ok2 {
			cartstr = append(cartstr, itemstr)
		}
		if err == nil && user != nil {
			err = sess.Set("cart", cartstr)
			if err == nil {
				c.Redirect("/cart")
			}
		}

	}
}

func Cart(c *ace.C) {
	w := c.Writer
	r := c.Request

	sess, err := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)

	user := sess.Get("email")
	cart := sess.Get("cart")
	if cartstring, ok := cart.([]string); ok {
		log.Println(cartstring)
		itemsincart := GetCart(cartstring)
		log.Println(itemsincart)
		max := len(itemsincart)

		if user != nil && err == nil {
			render := render.New(render.Options{})
			if max == 2 {
				render.HTML(c.Writer, http.StatusOK, "cart", map[string]interface{}{
					"item0": itemsincart[0],
					"item1": itemsincart[1],
				})
			} else if max == 1 {
				render.HTML(c.Writer, http.StatusOK, "cart", map[string]interface{}{
					"item0": itemsincart[0],
				})
			}

		}
	}
}

func GetCart(cart []string) []Item {
	var cartitems []Item

	for _, elem := range cart {
		cartitem := Item{}
		err := db.Get(&cartitem, "SELECT * FROM items WHERE id=$1", elem)
		if err != nil {
			log.Print("This must be the problem")
		}

		cartitems = append(cartitems, cartitem)
	}

	return cartitems
}

func Buy(c *ace.C) {
	stripe.SetKey("sk_test_s7zyOcXwo4E9YLlBXpL4Ie44")

}
