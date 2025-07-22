package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func ConnectDB(url string) error {
	db, err := sql.Open("sqlite3", url)
	if err != nil {
		return err
	}

	DB = db

	return nil
}
