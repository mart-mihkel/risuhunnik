package database

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

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

	cs, err := scanConundrums(rows)
	if err != nil {
		return nil, err
	}

	return cs, nil
}

func GetConundrumsByTag(t string) ([]Conundrum, error) {
	q := "SELECT * FROM conundrums WHERE instr(tags, ?) > 0"

	rows, err := DB.Query(q, t)
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

	row := DB.QueryRow(q, id)

	var text string
	var tags string
	var stars int
	var verified bool

	err := row.Scan(&id, &text, &tags, &stars, &verified)
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

func InsertConundrum(c *Conundrum) (int, error) {
	q := "INSERT INTO conundrums (text, tags, verified, stars) VALUES (?, ?, ?, ?)"

	res, err := DB.Exec(q, c.Text, strings.Join(c.Tags, " "), c.Verified, c.Stars)
	if err != nil {
		return 0, fmt.Errorf("failed to insert conundrum: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return int(id), nil
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
	var cs []Conundrum
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

		cs = append(cs, Conundrum{
			Id:       id,
			Text:     text,
			Tags:     strings.Split(tags, " "),
			Verified: verified,
			Stars:    stars,
		})
	}

	return cs, nil
}
