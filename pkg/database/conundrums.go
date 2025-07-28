package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Conundrum struct {
	Id       int
	Text     string
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

func GetConundrum(id int) (*Conundrum, error) {
	q := "SELECT * FROM conundrums WHERE id = ?"

	row := Db.QueryRow(q, id)

	var c Conundrum

	err := row.Scan(&c.Id, &c.Text, &c.Verified)
	if err != nil {
		return nil, fmt.Errorf("failed on scannig row: %w", err)
	}

	return &c, nil
}

func InsertConundrum(c *Conundrum) (int, error) {
	q := "INSERT INTO conundrums (text, verified) VALUES (?, ?)"

	res, err := Db.Exec(q, c.Text, c.Verified)
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
	q := "UPDATE conundrums SET text = ?, verified = ? WHERE id = ?"

	_, err := Db.Exec(q, c.Text, c.Verified, c.Id)
	if err != nil {
		return fmt.Errorf("failed to update conundrum: %w", err)
	}

	return nil
}

func scanConundrums(rows *sql.Rows) ([]Conundrum, error) {
	var cs []Conundrum
	for rows.Next() {
		var c Conundrum

		err := rows.Scan(&c.Id, &c.Text, &c.Verified)
		if err != nil {
			return nil, fmt.Errorf("failed on scannig row: %w", err)
		}

		cs = append(cs, c)
	}

	return cs, nil
}
