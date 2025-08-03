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
	Id          int
	ConundrumId int
	Comment     string
	Date        time.Time
}

func GetAllConundrums() ([]Conundrum, error) {
	q := "SELECT * FROM conundrums"

	rows, err := Db.Query(q)
	if err != nil {
		return nil, fmt.Errorf("couldn't get conundrums: %w", err)
	}

	defer rows.Close()

	conundrums, err := scanConundrums(rows)
	if err != nil {
		return nil, err
	}

	return conundrums, nil
}

func GetConundrum(id int) (*Conundrum, error) {
	q := "SELECT * FROM conundrums WHERE id = ?"

	var c Conundrum
	err := Db.QueryRow(q, id).Scan(
		&c.Id, &c.Text, &c.Verified,
		&c.Stars, &c.Date,
	)

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

	var comments []Comment
	for rows.Next() {
		var comment Comment
		err = rows.Scan(
			&comment.Id, &comment.ConundrumId,
			&comment.Comment, &comment.Date,
		)

		if err != nil {
			return nil, fmt.Errorf("failed on scannig row: %w", err)
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

func StarConundrum(id int) (*Conundrum, error) {
	q := "UPDATE conundrums SET stars = stars + 1 WHERE ID = ? RETURNING *"
	return starConundrum(id, q)
}

func UnStarConundrum(id int) (*Conundrum, error) {
	q := "UPDATE conundrums SET stars = stars - 1 WHERE ID = ? RETURNING *"
	return starConundrum(id, q)
}

func starConundrum(id int, q string) (*Conundrum, error) {
	var c Conundrum
	err := Db.QueryRow(q, id).Scan(
		&c.Id, &c.Text, &c.Verified,
		&c.Stars, &c.Date,
	)

	if err != nil {
		return nil, fmt.Errorf("failed on scannig row: %w", err)
	}

	return &c, nil
}

func InsertConundrum(text string) error {
	q := "INSERT INTO conundrums (text) VALUES (?)"

	_, err := Db.Exec(q, text)
	if err != nil {
		return fmt.Errorf("failed to insert conundrum: %w", err)
	}

	return nil
}

func InsertComment(conundrumId int, comment string) error {
	q := "INSERT INTO comments (cid, comment) VALUES (?, ?)"

	_, err := Db.Exec(q, conundrumId, comment)
	if err != nil {
		return fmt.Errorf("failed to insert comment: %w", err)
	}

	return nil
}

func scanConundrums(rows *sql.Rows) ([]Conundrum, error) {
	var conundrums []Conundrum
	for rows.Next() {
		var c Conundrum
		err := rows.Scan(&c.Id, &c.Text, &c.Verified, &c.Stars, &c.Date)
		if err != nil {
			return nil, fmt.Errorf("failed on scannig row: %w", err)
		}

		conundrums = append(conundrums, c)
	}

	return conundrums, nil
}
