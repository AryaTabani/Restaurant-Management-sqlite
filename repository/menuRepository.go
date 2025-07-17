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
	err := db.DB.QueryRow(`SELECT id FROM menus WHERE menu_id = ?`, menuID).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
