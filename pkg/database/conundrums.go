package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Conundrum struct {
	Id       int
	Text     string
	Verified bool
	Stars    int
	Date     time.Time
}

type Comment struct {
	Id      int
	Cid     int
	Comment string
	Date    time.Time
}

func GetAllConundrums() ([]Conundrum, error) {
	q := "SELECT * FROM conundrums"

	rows, err := Db.Query(q)
	if err != nil {
		return nil, fmt.Errorf("couldn't get conundrums: %w", err)
	}

	defer rows.Close()

	cs, err := scanConundrums(rows)
	if err != nil {
		return nil, err
	}

	return cs, nil
}

func GetConundrum(id int) (*Conundrum, error) {
	q := "SELECT * FROM conundrums WHERE id = ?"

	row := Db.QueryRow(q, id)

	var c Conundrum

	err := row.Scan(&c.Id, &c.Text, &c.Verified, &c.Stars, &c.Date)
	if err != nil {
		return nil, fmt.Errorf("failed on scannig row: %w", err)
	}

	return &c, nil
}

func GetConundrumComments(id int) ([]Comment, error) {
	q := "SELECT * FROM comments WHERE cid = ?"

	rows, err := Db.Query(q, id)
	if err != nil {
		return nil, fmt.Errorf("couldn't get comments: %w", err)
	}

	defer rows.Close()

	var cs []Comment
	for rows.Next() {
		var c Comment

		err = rows.Scan(&c.Id, &c.Cid, &c.Comment, &c.Date)
		if err != nil {
			return nil, fmt.Errorf("failed on scannig row: %w", err)
		}

		cs = append(cs, c)
	}

	return cs, nil
}

func StarConundrum(id int) (*Conundrum, error) {
	q := "UPDATE conundrums SET stars = stars + 1 WHERE ID = ? RETURNING *"

	row := Db.QueryRow(q, id)

	var c Conundrum

	err := row.Scan(&c.Id, &c.Text, &c.Verified, &c.Stars, &c.Date)
	if err != nil {
		return nil, fmt.Errorf("failed on scannig row: %w", err)
	}

	return &c, nil
}

func InsertConundrum(t string) error {
	q := "INSERT INTO conundrums (text) VALUES (?)"

	_, err := Db.Exec(q, t)
	if err != nil {
		return fmt.Errorf("failed to insert conundrum: %w", err)
	}

	return nil
}

func InsertComment(cid int, co string) error {
	q := "INSERT INTO comments (cid, comment) VALUES (?, ?)"

	_, err := Db.Exec(q, cid, co)
	if err != nil {
		return fmt.Errorf("failed to insert comment: %w", err)
	}

	return nil
}

func scanConundrums(rows *sql.Rows) ([]Conundrum, error) {
	var cs []Conundrum
	for rows.Next() {
		var c Conundrum

		err := rows.Scan(&c.Id, &c.Text, &c.Verified, &c.Stars, &c.Date)
		if err != nil {
			return nil, fmt.Errorf("failed on scannig row: %w", err)
		}

		cs = append(cs, c)
	}

	return cs, nil
}
