package database

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func CheckToken(token string) (bool, error) {
	q := "SELECT EXISTS (SELECT * FROM tokens WHERE token = ?)"

	var valid bool

	err := Db.QueryRow(q, token).Scan(&valid)
	if err != nil {
		return false, fmt.Errorf("getting token valid: %w", err)
	}

	return valid, nil
}
