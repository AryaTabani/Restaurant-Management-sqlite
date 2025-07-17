package repository

import (
	"database/sql"
	"errors"

	db "example.com/m/v2/DB"
)

var (
	ErrMenuNotFound = errors.New("the specified menu does not exist")
)

func MenuExists(menuID string) (bool, error) {
	var id int64
	// This query checks for existence without fetching all the data.
	err := db.DB.QueryRow(`SELECT id FROM menus WHERE menu_id = ?`, menuID).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // Menu does not exist, which is not a database error.
		}
		return false, err // A different database error occurred.
	}
	return true, nil // Menu exists.
}
