package main

import (
	"fmt"
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
			} else {
				badstr := "You need to be signed in to add to cart"
				BadPage(c, badstr)
			}
		}

	} else {
		badstr := "You need to be signed in to add to cart"
		BadPage(c, badstr)
	}
}

func DeleteFromCart(c *ace.C) {
	w := c.Writer
	r := c.Request

	sess, err := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)
	log.Println(sess)
	deleteid := r.FormValue("ajax_post_data")
	log.Println(deleteid)

	cartcookie := sess.Get("cart")
	log.Println(cartcookie)
	if a, ok := cartcookie.([]string); ok {
		log.Println(a)
		var i int
		if _, err := fmt.Sscanf(deleteid, "%5d", &i); err == nil {
			aints := Stringtoint(a)
			for _, elem := range aints {
				if elem == i {
					elempos := aints.pos(elem)
					a[elempos] = a[len(a)-1]
					a = a[:len(a)-1]
				}
			}
		}
		log.Println(a)
		cart := GetCart(a)
		log.Println(cart)
		err = sess.Set("cart", a)

		RenderCart(c, cart)

		log.Println(err)

	}
}

func Cart(c *ace.C) {
	w := c.Writer
	r := c.Request

	sess, err := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)

	user := sess.Get("email")
	cart := sess.Get("cart")

	itemsincart := GetCart(cart)
	log.Println(itemsincart)
	max := len(itemsincart)

	if max > 5 {
		if cartstr, ok := cart.([]string); ok {

			newcart := cartstr[0:4]
			sess.Delete("cart")
			sess.Set("cart", newcart)
		}
	}

	if user != nil && err == nil {
		RenderCart(c, itemsincart)
	} else {
		badstr := "Sorry, you must be signed in to access this page"
		BadPage(c, badstr)
	}

}

func GetCart(cart interface{}) []Item {
	var cartitems []Item

	if cartstring, ok := cart.([]string); ok {
		for _, elem := range cartstring {
			cartitem := Item{}
			err := db.Get(&cartitem, "SELECT * FROM items WHERE id=$1", elem)
			if err != nil {
				log.Print("This must be the problem")
			}

			cartitems = append(cartitems, cartitem)
		}
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

func RenderCart(c *ace.C, itemsincart []Item) {
	max := len(itemsincart)
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
		SimplePage(c.Writer, c.Request, "cart")
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

type intSlice []int

func Stringtoint(strs []string) intSlice {
	var ints intSlice
	for _, elem := range strs {
		var i int

		if _, err := fmt.Sscanf(elem, "%5d", &i); err == nil {
			ints = append(ints, i)
		}

	}
	return ints
}

func (slice intSlice) pos(value int) int {
	for p, v := range slice {
		if v == value {
			return p
		}
	}
	return -1
}

func Checkout(c *ace.C) {
	r := c.Request
	w := c.Writer
	stripe.SetKey("sk_test_s7zyOcXwo4E9YLlBXpL4Ie44")
	r.ParseMultipartForm(5)

	sess, err := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)

	user := sess.Get("email")
	cartcookie := sess.Get("cart")
	badstr := "There was an error with your purchase."

	cart := GetCart(cartcookie)
	total := CartTotal(cart)
	totalcent := total * 100
	log.Println(totalcent)
	totalcentint64 := int64(totalcent)

	name := r.PostFormValue("nameoncard")
	cardnum := r.PostFormValue("cardnumber")
	expmonth := r.PostFormValue("expirationmonth")
	expyear := r.PostFormValue("expirationyear")

	if user != nil && err == nil {
		err = sess.Delete("cart")

		expmonthnum, _ := strconv.ParseInt(expmonth, 0, 64)
		expyearnum, _ := strconv.ParseInt(expyear, 0, 64)
		expyearnumint := int(expyearnum)
		expmonthnumint := int(expmonthnum)
		log.Println(totalcentint64)
		params := stripe.ChargeParams{
			Desc:     "Cart",
			Amount:   totalcentint64,
			Currency: "usd",
			Card: &stripe.CardParams{
				Name:     name,
				Number:   cardnum,
				ExpYear:  expyearnumint,
				ExpMonth: expmonthnumint,
			},
		}
		log.Println(params)
		charge, err := stripe.Charges.Create(&params)

		log.Println(charge)
		log.Println(err)
		if err == nil {
			render := render.New(render.Options{})
			render.HTML(c.Writer, http.StatusOK, "checkout", nil)
		}
	} else {
		BadPage(c, badstr)
	}

}

func Buy(c *ace.C) {
	var w = c.Writer
	var r = c.Request

	sess, err := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)

	user := sess.Get("email")
	cart := sess.Get("cart")

	itemsincart := GetCart(cart)
	total := CartTotal(itemsincart)
	badstr := "you must be signed in to access this page"

	if user != nil && err == nil {
		render := render.New(render.Options{})
		render.HTML(c.Writer, http.StatusOK, "buy", total)
	} else {
		BadPage(c, badstr)
	}

}
