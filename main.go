package main

import (
	"database/sql"
	"log"
	"net/http"
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
		var t string

		err := rows.Scan(&j.Joke, &t)
		if err != nil {
			log.Fatal(err)
		}

		j.Tags = strings.Split(t, ",")
		jokes = append(jokes, j)
	}

	return jokes
}

// func getTags(db *sql.DB) []string {
// 	tags := []string{}
// 	return tags
// }

func main() {
	db, err := sql.Open("sqlite3", "risuhunnik.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/", http.StripPrefix("/", fs))

	http.HandleFunc("/api/jokes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			tmpl := template.Must(template.ParseFiles("templates/joke.html"))
			jokes := getJokes(db)
			tmpl.Execute(w, jokes)
		}

		if r.Method == http.MethodPost {
			// TODO: post joke
		}
	})

	http.ListenAndServe(":8080", nil)
}
