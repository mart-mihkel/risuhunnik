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
	Id       int
	Joke     string
	Tags     []string
	Verified bool
	Stars    string
}

func getJokes(db *sql.DB) []joke {
	rows, err := db.Query("SELECT id, joke, tags, verified, stars FROM jokes")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var jokes []joke
	for rows.Next() {
		var j joke
		var tagstring string

		err := rows.Scan(&j.Id, &j.Joke, &tagstring, &j.Verified, &j.Stars)
		if err != nil {
			log.Fatal(err)
		}

		j.Tags = strings.Split(tagstring, ",")
		jokes = append(jokes, j)
	}

	return jokes
}

func postJoke(db *sql.DB, jokestring *string, tagstring *string) {
	_, err := db.Exec("INSERT INTO jokes (joke, tags) VALUES (?, ?)", *jokestring, *tagstring)
	// TODO: bubble error, unique constraint
	if err != nil {
		log.Fatal(err)
	}
}

func starJoke(db *sql.DB, id int64) {
	_, err := db.Exec("UPDATE jokes SET stars = stars + 1 WHERE id = ?", id)
	// TODO: bubble error
	if err != nil {
		log.Fatal(err)
	}
}

func jokesSearch(jokes []joke, search *string) []joke {
	var res []joke
	for _, j := range jokes {
		if strings.Contains(j.Joke, *search) {
			res = append(res, j)
		}
	}
	return res
}

func jokesTags(jokes []joke, tag *string) []joke {
	var res []joke
	for _, j := range jokes {
		if slices.Contains(j.Tags, *tag) {
			res = append(res, j)
		}
	}
	return res
}

func getTagCounts(db *sql.DB, counts *map[string]int) {
	rows, err := db.Query("SELECT tags FROM jokes")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var t string

		err := rows.Scan(&t)
		if err != nil {
			log.Fatal(err)
		}

		for tag := range strings.SplitSeq(t, ",") {
			(*counts)[tag]++
		}
	}
}

func main() {
	db, err := sql.Open("sqlite3", "risuhunnik.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fs := http.FileServer(http.Dir("../web/static/"))
	http.Handle("/", http.StripPrefix("/", fs))

	http.HandleFunc("/jokes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			jokes := getJokes(db)

			q := r.URL.Query()
			search := q.Get("search")
			if search != "" {
				jokes = jokesSearch(jokes, &search)
			}

			tag := q.Get("tag")
			if tag != "" {
				jokes = jokesTags(jokes, &tag)
			}

			tmpl := template.Must(template.ParseFiles("../web/templates/jokes.html"))
			tmpl.Execute(w, jokes)
		}

		if r.Method == http.MethodPost {
			q := r.URL.Query()
			likestr := q.Get("star")
			if likestr != "" {
				id, err := strconv.ParseInt(likestr, 10, 64)
				if err != nil {
					log.Printf("Recieved malformed url paramater star=%s, should be int\n", likestr)
					return
				}

				starJoke(db, id)

				// TODO: return new number of likes
				return
			}

			jokestring := r.FormValue("joke")
			tagstring := r.FormValue("tags")
			postJoke(db, &jokestring, &tagstring)
		}
	})

	http.HandleFunc("/tags", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			return
		}

		tags := make(map[string]int)
		getTagCounts(db, &tags)

		tmpl := template.Must(template.ParseFiles("../web/templates/tags.html"))
		tmpl.Execute(w, tags)
	})

	http.HandleFunc("/search-modal", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			return
		}

		openstr := r.URL.Query().Get("open")
		open, err := strconv.ParseBool(openstr)
		if err != nil {
			log.Printf("Recieved malformed url paramater open=%s, should be boolean\n", openstr)
			return
		}

		tmpl := template.Must(template.ParseFiles("../web/templates/search-modal.html"))
		tmpl.Execute(w, open)
	})

	http.HandleFunc("/add-modal", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			return
		}

		openstr := r.URL.Query().Get("open")
		open, err := strconv.ParseBool(openstr)
		if err != nil {
			log.Printf("Recieved malformed url paramater open=%s, should be boolean\n", openstr)
			return
		}

		tmpl := template.Must(template.ParseFiles("../web/templates/add-modal.html"))
		tmpl.Execute(w, open)
	})

	http.ListenAndServe(":8080", nil)
}
