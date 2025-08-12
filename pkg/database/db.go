package database

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Conundrum struct {
	Id       int
	Text     string
	Author   string
	Verified bool
	Stars    int
	Date     time.Time
}

type Comment struct {
	Id          int
	ConundrumId int
	Comment     string
	Author      string
	Date        time.Time
}

var Db *sql.DB

func ConnectDB(url string) error {
	db, err := sql.Open("sqlite3", url)
	if err != nil {
		return err
	}

	Db = db

	return nil
}
