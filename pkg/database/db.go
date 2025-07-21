package database

import (
	"database/sql"
	"fmt"
	"strings"

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

type Conundrum struct {
	Id       int
	Text     string
	Tags     []string
	Stars    int
	Verified bool
}

func GetAllConundrums() ([]Conundrum, error) {
	q := "SELECT * FROM conundrums"

	rows, err := DB.Query(q)
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

func GetConundrumsByTag(tag string) ([]Conundrum, error) {
	q := "SELECT * FROM conundrums WHERE instr(tags, ?) > 0"

	rows, err := DB.Query(q, tag)
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

func GetConundrumsBySubstring(substr string) ([]Conundrum, error) {
	q := "SELECT * FROM conundrums WHERE instr(text, ?) > 0"

	rows, err := DB.Query(q, substr)
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

	row := DB.QueryRow(q, id)

	var text string
	var tags string
	var verified bool
	var stars int

	err := row.Scan(&id, &text, &tags, &verified, &stars)
	if err != nil {
		return nil, fmt.Errorf("failed on scannig row: %w", err)
	}

	return &Conundrum{
		Id:       id,
		Text:     text,
		Tags:     strings.Split(tags, " "),
		Verified: verified,
		Stars:    stars,
	}, nil
}

func InsertConundrum(c *Conundrum) error {
	q := "INSERT INTO conundrums (text, tags, verified, stars) VALUES (?, ?, ?, ?)"

	_, err := DB.Exec(q, c.Text, strings.Join(c.Tags, " "), c.Verified, c.Stars)
	if err != nil {
		return fmt.Errorf("failed to insert conundrum: %w", err)
	}

	return nil
}

func UpdateConundrum(c *Conundrum) error {
	q := "UPDATE conundrums SET text = ?, tags = ?, verified = ?, stars = ? WHERE id = ?"

	_, err := DB.Exec(q, c.Text, strings.Join(c.Tags, " "), c.Verified, c.Stars, c.Id)
	if err != nil {
		return fmt.Errorf("failed to update conundrum: %w", err)
	}

	return nil
}

func scanConundrums(rows *sql.Rows) ([]Conundrum, error) {
	var conundrums []Conundrum
	for rows.Next() {
		var id int
		var text string
		var tags string
		var stars int
		var verified bool

		err := rows.Scan(&id, &text, &tags, &stars, &verified)
		if err != nil {
			return nil, fmt.Errorf("failed on scannig row: %w", err)
		}

		conundrums = append(conundrums, Conundrum{
			Id:       id,
			Text:     text,
			Tags:     strings.Split(tags, " "),
			Verified: verified,
			Stars:    stars,
		})
	}

	return conundrums, nil
}
