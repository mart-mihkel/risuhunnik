package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

func ConnectDB(url string) error {
	db, err := sql.Open("sqlite3", url)
	if err != nil {
		return err
	}

	Db = db

	return nil
}
