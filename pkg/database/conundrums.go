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

func GetConundrumsByTag(t string) ([]Conundrum, error) {
	q := "SELECT * FROM conundrums WHERE instr(tags, ?) > 0"

	rows, err := Db.Query(q, t)
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
	var t string

	err := row.Scan(&c.Id, &c.Text, &t, &c.Stars, &c.Verified)
	if err != nil {
		return nil, fmt.Errorf("failed on scannig row: %w", err)
	}

	c.Tags = strings.Split(t, " ")

	return &c, nil
}

func InsertConundrum(c *Conundrum) (int, error) {
	q := "INSERT INTO conundrums (text, tags, verified, stars) VALUES (?, ?, ?, ?)"
	t := strings.Join(c.Tags, " ")

	res, err := Db.Exec(q, c.Text, t, c.Verified, c.Stars)
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
	q := "UPDATE conundrums SET text = ?, tags = ?, stars = ?, verified = ? WHERE id = ?"
	t := strings.Join(c.Tags, " ")

	_, err := Db.Exec(q, c.Text, t, c.Stars, c.Verified, c.Id)
	if err != nil {
		return fmt.Errorf("failed to update conundrum: %w", err)
	}

	return nil
}

func StarConundrum(id int) (*Conundrum, error) {
	q := "UPDATE conundrums SET stars = stars + 1 WHERE id = ? RETURNING *"

	var c Conundrum
	var t string

	row := Db.QueryRow(q, id)
	err := row.Scan(&c.Id, &c.Text, &t, &c.Stars, &c.Verified)
	if err != nil {
		return nil, fmt.Errorf("failed to star conundrum: %w", err)
	}

	c.Tags = strings.Split(t, " ")

	return &c, nil
}

func scanConundrums(rows *sql.Rows) ([]Conundrum, error) {
	var cs []Conundrum
	for rows.Next() {
		var c Conundrum
		var t string

		err := rows.Scan(&c.Id, &c.Text, &t, &c.Stars, &c.Verified)
		if err != nil {
			return nil, fmt.Errorf("failed on scannig row: %w", err)
		}

		c.Tags = strings.Split(t, " ")
		cs = append(cs, c)
	}

	return cs, nil
}
