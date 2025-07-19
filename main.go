package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

type joke struct {
	Joke string
}

func getJokes(db *sql.DB) []joke {
	rows, err := db.Query("SELECT joke FROM jokes")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var jokes []joke
	for rows.Next() {
		var j joke

		err := rows.Scan(&j.Joke)
		if err != nil {
			log.Fatal(err)
		}

		jokes = append(jokes, j)
	}

	return jokes
}

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
