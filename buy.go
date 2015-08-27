package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

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
			if max == 5 {
				total := CartTotal(itemsincart)
				render.HTML(c.Writer, http.StatusOK, "cart", map[string]interface{}{
					"total": total,
					"item0": itemsincart[0],
					"item1": itemsincart[1],
					"item2": itemsincart[2],
					"item3": itemsincart[3],
					"item4": itemsincart[4],
				})
			} else if max == 4 {
				total := CartTotal(itemsincart)
				render.HTML(c.Writer, http.StatusOK, "cart", map[string]interface{}{
					"total": total,
					"item0": itemsincart[0],
					"item1": itemsincart[1],
					"item2": itemsincart[2],
					"item3": itemsincart[3],
				})
			} else if max == 3 {
				total := CartTotal(itemsincart)
				render.HTML(c.Writer, http.StatusOK, "cart", map[string]interface{}{
					"total": total,
					"item0": itemsincart[0],
					"item1": itemsincart[1],
					"item2": itemsincart[2],
				})
			} else if max == 2 {
				total := CartTotal(itemsincart)
				render.HTML(c.Writer, http.StatusOK, "cart", map[string]interface{}{
					"total": total,
					"item0": itemsincart[0],
					"item1": itemsincart[1],
				})
			} else if max == 1 {
				total := CartTotal(itemsincart)
				render.HTML(c.Writer, http.StatusOK, "cart", map[string]interface{}{
					"total": total,
					"item0": itemsincart[0],
				})
			} else if max == 0 {
				SimplePage(w, r, "cart")
			} else if max > 5 {
				total := CartTotal(itemsincart)
				maxitems := true
				render.HTML(c.Writer, http.StatusOK, "cart", map[string]interface{}{
					"total": total,
					"max":   maxitems,
					"item0": itemsincart[0],
					"item1": itemsincart[1],
					"item2": itemsincart[2],
					"item3": itemsincart[3],
					"item4": itemsincart[4],
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

func CartTotal(cart []Item) float64 {
	total := 0.0

	for _, elem := range cart {
		price := strings.Trim(elem.Price, "$")
		if n, err := strconv.ParseFloat(price, 64); err == nil {
			total += n
		}
	}

	return total
}

func Buy(c *ace.C) {
	stripe.SetKey("sk_test_s7zyOcXwo4E9YLlBXpL4Ie44")

}
