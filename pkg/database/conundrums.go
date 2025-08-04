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

func GetAuthorConundrums(author string) ([]Conundrum, error) {
	q := "SELECT * FROM conundrums WHERE author = ?"

	rows, err := Db.Query(q, author)
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

func GetAuthorStars(author string) (int, error) {
	q := "SELECT SUM(stars) FROM conundrums WHERE author = ?"

	var stars int
	err := Db.QueryRow(q, author).Scan(&stars)
	if err != nil {
		return 0, fmt.Errorf("couldn't get conundrums: %w", err)
	}

	return stars, nil
}

func GetConundrum(id int) (*Conundrum, error) {
	q := "SELECT * FROM conundrums WHERE id = ?"

	var c Conundrum
	err := Db.QueryRow(q, id).Scan(
		&c.Id, &c.Text, &c.Author,
		&c.Verified, &c.Stars, &c.Date,
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
			&comment.Id, &comment.ConundrumId, &comment.Comment,
			&comment.Author, &comment.Date,
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
		&c.Id, &c.Text, &c.Author,
		&c.Verified, &c.Stars, &c.Date,
	)

	if err != nil {
		return nil, fmt.Errorf("failed on scannig row: %w", err)
	}

	return &c, nil
}

func InsertConundrum(text string, author string) error {
	q := "INSERT INTO conundrums (text, author) VALUES (?, ?)"

	_, err := Db.Exec(q, text, author)
	if err != nil {
		return fmt.Errorf("failed to insert conundrum: %w", err)
	}

	return nil
}

func InsertComment(conundrumId int, comment string, author string) error {
	q := "INSERT INTO comments (cid, comment, author) VALUES (?, ?, ?)"

	_, err := Db.Exec(q, conundrumId, comment, author)
	if err != nil {
		return fmt.Errorf("failed to insert comment: %w", err)
	}

	return nil
}

func RandomAuthor() (string, error) {
	q1 := "SELECT author FROM authors ORDER BY RANDOM() LIMIT 1"

	var author string
	err := Db.QueryRow(q1).Scan(&author)
	if err != nil {
		return "", fmt.Errorf("failed on scanning row: %w", err)
	}

	return author, nil
}

func scanConundrums(rows *sql.Rows) ([]Conundrum, error) {
	var conundrums []Conundrum
	for rows.Next() {
		var c Conundrum
		err := rows.Scan(
			&c.Id, &c.Text, &c.Author,
			&c.Verified, &c.Stars, &c.Date,
		)

		if err != nil {
			return nil, fmt.Errorf("failed on scannig row: %w", err)
		}

		conundrums = append(conundrums, c)
	}

	return conundrums, nil
}
