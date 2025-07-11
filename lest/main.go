package main

import (
    "database/sql"
    "html/template"
    "log"
    "net/http"
    "strings"
    _ "github.com/mattn/go-sqlite3"
)

type joke struct {
    Joke string
    Tags []string
}

func main() {
    db, err := sql.Open("sqlite3", "./lest.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    query := "SELECT joke, tags FROM jokes"
    rows, err := db.Query(query)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    var jokes []joke
    for rows.Next() {
        var j joke
        var strtags string
        err := rows.Scan(&j.Joke, &strtags)
        if err != nil {
            log.Fatal(err)
        }
        j.Tags = strings.Split(strtags, ",")
        jokes = append(jokes, j)
    }

    tmpl := template.Must(template.New("joke").Parse(`
        {{range .}}
          <div>
            <p>{{.Joke}}</p>
            <ul>
              {{range .Tags}}
                <li>{{.}}</li>
              {{end}}
            </ul>
          </div>
        {{end}}`))

    fs := http.FileServer(http.Dir("static/"))
    http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
        tmpl.Execute(w, jokes)
    })

    http.Handle("/static/", http.StripPrefix("/static/", fs))

    http.ListenAndServe(":8080", nil)
}
