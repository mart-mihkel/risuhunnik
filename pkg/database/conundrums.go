package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func GetAllConundrums() ([]Conundrum, error) {
	q := "SELECT * FROM conundrums"

	rows, err := Db.Query(q)
	if err != nil {
		return nil, fmt.Errorf("getting conundrums: %w", err)
	}

	defer rows.Close()

	conundrums, err := scanConundrums(rows)
	if err != nil {
		return nil, err
	}

	return conundrums, nil
}

func GetVerifiedConundrums() ([]Conundrum, error) {
	q := "SELECT * FROM conundrums WHERE verified = 1"

	rows, err := Db.Query(q)
	if err != nil {
		return nil, fmt.Errorf("getting conundrums: %w", err)
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
		&c.Id, &c.Text, &c.Author,
		&c.Verified, &c.Stars, &c.Date,
	)

	if err != nil {
		return nil, fmt.Errorf("scannig row: %w", err)
	}

	return &c, nil
}

func GetConundrumCount() (int, error) {
	q := "SELECT COUNT(*) FROM conundrums"

	var c int
	err := Db.QueryRow(q).Scan(&c)
	if err != nil {
		return 0, fmt.Errorf("scannig row: %w", err)
	}

	return c, nil
}

func GetConundrumComments(id int) ([]Comment, error) {
	q := "SELECT * FROM comments WHERE cid = ?"

	rows, err := Db.Query(q, id)
	if err != nil {
		return nil, fmt.Errorf("getting comments: %w", err)
	}

	defer rows.Close()

	comments, err := scanComments(rows)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func ToggleVerifyConundrum(id int) (*Conundrum, error) {
	q := `UPDATE conundrums
		SET verified = 1 - verified
		WHERE ID = ?
		RETURNING *`

	return updateConundrumById(id, q)
}

func StarConundrum(id int) (*Conundrum, error) {
	q := "UPDATE conundrums SET stars = stars + 1 WHERE ID = ? RETURNING *"
	return updateConundrumById(id, q)
}

func UnStarConundrum(id int) (*Conundrum, error) {
	q := "UPDATE conundrums SET stars = stars - 1 WHERE ID = ? RETURNING *"
	return updateConundrumById(id, q)
}

func updateConundrumById(id int, q string) (*Conundrum, error) {
	var c Conundrum
	err := Db.QueryRow(q, id).Scan(
		&c.Id, &c.Text, &c.Author,
		&c.Verified, &c.Stars, &c.Date,
	)

	if err != nil {
		return nil, fmt.Errorf("scannig row: %w", err)
	}

	return &c, nil
}

func InsertConundrum(text string, author string) error {
	q := "INSERT INTO conundrums (text, author) VALUES (?, ?)"

	_, err := Db.Exec(q, text, author)
	if err != nil {
		return fmt.Errorf("inserting conundrum: %w", err)
	}

	return nil
}

func InsertComment(conundrumId int, comment string, author string) error {
	q := "INSERT INTO comments (cid, comment, author) VALUES (?, ?, ?)"

	_, err := Db.Exec(q, conundrumId, comment, author)
	if err != nil {
		return fmt.Errorf("inserting comment: %w", err)
	}

	return nil
}

func scanConundrums(rows *sql.Rows) ([]Conundrum, error) {
	var conundrums []Conundrum
	for rows.Next() {
		var c Conundrum
		err := rows.Scan(
			&c.Id, &c.Text, &c.Author, &c.Verified, &c.Stars, &c.Date,
		)

		if err != nil {
			return nil, fmt.Errorf("scannig row: %w", err)
		}

		conundrums = append(conundrums, c)
	}

	return conundrums, nil
}

func scanComments(rows *sql.Rows) ([]Comment, error) {
	var comments []Comment
	for rows.Next() {
		var c Comment
		err := rows.Scan(
			&c.Id, &c.ConundrumId, &c.Comment, &c.Author, &c.Date,
		)

		if err != nil {
			return nil, fmt.Errorf("scannig row: %w", err)
		}

		comments = append(comments, c)
	}

	return comments, nil
}
