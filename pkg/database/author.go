package database

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func RandomAuthor() (string, error) {
	q1 := "SELECT author FROM authors ORDER BY RANDOM() LIMIT 1"

	var author string
	err := Db.QueryRow(q1).Scan(&author)
	if err != nil {
		return "", fmt.Errorf("failed on scanning row: %w", err)
	}

	return author, nil
}

func GetAuthorConundrums(author string) ([]Conundrum, error) {
	q := "SELECT * FROM conundrums WHERE author = ?"

	rows, err := Db.Query(q, author)
	if err != nil {
		return nil, fmt.Errorf("couldn't get author conundrums: %w", err)
	}

	defer rows.Close()

	conundrums, err := scanConundrums(rows)
	if err != nil {
		return nil, err
	}

	return conundrums, nil
}

func GetAuthorComments(author string) ([]Comment, error) {
	q := "SELECT * FROM comments WHERE author = ?"

	rows, err := Db.Query(q, author)
	if err != nil {
		return nil, fmt.Errorf("couldn't get author comments: %w", err)
	}

	defer rows.Close()

	comments, err := scanComments(rows)
	if err != nil {
		return nil, err
	}

	return comments, nil
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
