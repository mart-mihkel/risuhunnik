package main

import (
	"database/sql"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

type joke struct {
	Joke string
	Tags []string
}

func getJokes(db *sql.DB) []joke {
	rows, err := db.Query("SELECT joke, tags FROM jokes")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var jokes []joke
	for rows.Next() {
		var j joke
		var tagstring string

		err := rows.Scan(&j.Joke, &tagstring)
		if err != nil {
			log.Fatal(err)
		}

		j.Tags = strings.Split(tagstring, ",")
		jokes = append(jokes, j)
	}

	return jokes
}

func jokesSearch(jokes []joke, search string) []joke {
	var res []joke
	for _, j := range jokes {
		if strings.Contains(j.Joke, search) {
			res = append(res, j)
		}
	}
	return res
}

func jokesTags(jokes []joke, tag string) []joke {
	var res []joke
	for _, j := range jokes {
		if slices.Contains(j.Tags, tag) {
			res = append(res, j)
		}
	}
	return res
}

func getTags(db *sql.DB) map[string]int {
	rows, err := db.Query("SELECT tags FROM jokes")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var t string

		err := rows.Scan(&t)
		if err != nil {
			log.Fatal(err)
		}

		for tag := range strings.SplitSeq(t, ",") {
			counts[tag]++
		}
	}

	return counts
}

func main() {
	db, err := sql.Open("sqlite3", "risuhunnik.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fs := http.FileServer(http.Dir("../web/static/"))
	http.Handle("/", http.StripPrefix("/", fs))

	http.HandleFunc("/api/jokes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			return
		}

		jokes := getJokes(db)

		q := r.URL.Query()
		search := q.Get("search")
		if search != "" {
			jokes = jokesSearch(jokes, search)
		}

		tag := q.Get("tag")
		if tag != "" {
			jokes = jokesTags(jokes, tag)
		}

		tmpl := template.Must(template.ParseFiles("../web/templates/jokes.html"))
		tmpl.Execute(w, jokes)
	})

	http.HandleFunc("/api/tags", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			return
		}

		tags := getTags(db)

		tmpl := template.Must(template.ParseFiles("../web/templates/tags.html"))
		tmpl.Execute(w, tags)
	})

	http.HandleFunc("/api/modal", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			return
		}

		openstr := r.URL.Query().Get("open")
		open, err := strconv.ParseBool(openstr)
		if err != nil {
			log.Printf("Recieved malformed url paramater open=%s, should be boolean\n", openstr)
			return
		}

		tmpl := template.Must(template.ParseFiles("../web/templates/modal.html"))
		tmpl.Execute(w, open)
	})

	http.ListenAndServe(":8080", nil)
}
