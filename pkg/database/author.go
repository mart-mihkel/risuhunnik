package database

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func RandomAuthor() (string, error) {
	q := "SELECT author FROM authors ORDER BY RANDOM() LIMIT 1"

	var author string
	err := Db.QueryRow(q).Scan(&author)
	if err != nil {
		return "", fmt.Errorf("scanning row: %w", err)
	}

	return author, nil
}

func GetAuthorConundrums(author string) ([]Conundrum, error) {
	q := "SELECT * FROM conundrums WHERE author = ?"

	rows, err := Db.Query(q, author)
	if err != nil {
		return nil, fmt.Errorf("getting author conundrums: %w", err)
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
		return nil, fmt.Errorf("getting author comments: %w", err)
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
		return 0, fmt.Errorf("getting author stars: %w", err)
	}

	return stars, nil
}
