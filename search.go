package main

import (
	"log"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
	"github.com/plimble/ace"
	"github.com/unrolled/render"
)

func Substring(s string, substring string) bool {
	lowers := strings.ToLower(s)
	lowersubstring := strings.ToLower(substring)
	compare := strings.Contains(lowers, lowersubstring)
	tf := false
	if compare {
		tf = true
	}
	return tf

}

func Compare(arr []Item, query string) (bool, []Item) {
	var tf bool
	results := []Item{}
	var sub bool
	var str bool

	for _, elem := range arr {
		sub = Substring(elem.Title, query)
		str = strings.EqualFold(elem.Title, query)
		log.Println(sub, str)
		if sub || str {
			tf = true
			results = append(results, elem)
		}
	}

	return tf, results
}

func Search(c *ace.C) {
	r := c.Request
	usersearch := r.FormValue("usersearch")

	itemsearch := []Item{}
	err := db.Select(&itemsearch, "SELECT * FROM items ORDER BY id ASC")
	if err != nil {
		log.Println("Is this the problem?")
		log.Println(err)
		return
	}

	compared, results := Compare(itemsearch, usersearch)
	log.Println("%#v\n", results)
	max := len(results)

	if compared {
		render := render.New(render.Options{})
		if max == 3 {
			render.HTML(c.Writer, http.StatusOK, "results", map[string]interface{}{
				"result0": results[0],
				"result1": results[1],
				"result2": results[2],
			})
		} else if max == 2 {
			render.HTML(c.Writer, http.StatusOK, "results", map[string]interface{}{
				"result0": results[0],
				"result1": results[1],
			})
		} else if max == 1 {
			render.HTML(c.Writer, http.StatusOK, "results", map[string]interface{}{
				"result0": results[0],
			})
		}
	} else {
		c.Redirect("/noresults")
	}
}
