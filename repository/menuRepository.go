package repository

import (
	"database/sql"
	"errors"

	db "example.com/m/v2/DB"
	"example.com/m/v2/models"
)

var (
	ErrMenuNotFound = errors.New("the specified menu does not exist")
	ErrMenuExist    = errors.New("the specified menu does not exist")
)

func MenuExists(menuID string) (bool, error) {
	var id int64
	err := db.DB.QueryRow(`SELECT id FROM menus WHERE menuid = ?`, menuID).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
func CreateMenu(menu *models.Menu) (int64, error) {
	stmt, err := db.DB.Prepare(`
		INSERT INTO menus (name, category, startdate, enddate, createdat, updatedat, menuid)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		menu.Name,
		menu.Category,
		menu.Start_Date,
		menu.End_Date,
		menu.Created_at,
		menu.Updated_at,
		menu.Menu_id,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}
func GetMenu(menuID string) (*models.Menu, error) {
	var menu models.Menu
	err := db.DB.QueryRow(`
		SELECT id, name, category, startdate, enddate, createdat, updatedat, menuid
		FROM menus WHERE menu_id = ?`, menuID).Scan(
		&menu.ID, &menu.Name, &menu.Category, &menu.Start_Date, &menu.End_Date,
		&menu.Created_at, &menu.Updated_at, &menu.Menu_id,
	)
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

func UpdateMenu(menu *models.Menu) error {
	stmt, err := db.DB.Prepare(`
		UPDATE menus SET name = ?, category = ?, startdate = ?, enddate = ?, updatedat = ?
		WHERE id = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		menu.Name, menu.Category, menu.Start_Date, menu.End_Date, menu.Updated_at,
		menu.ID,
	)
	return err
}
