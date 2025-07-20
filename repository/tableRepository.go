package repository

import (
	"database/sql"

	db "example.com/m/v2/DB"
	"example.com/m/v2/models"
)

func CreateTable(table *models.Table) (int64, error) {
	stmt, err := db.DB.Prepare(`
		INSERT INTO tables (numberofguests, tablenumber, createdat, updatedat, tableid)
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		table.Number_of_guests, table.Table_number, table.Created_at, table.Updated_at, table.Table_id,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}
func GetTableById(tableId string) (*models.Table, error) {
	var t models.Table
	row := db.DB.QueryRow(`
		SELECT id,numberofguests, tablenumber, createdat, updatedat, tableid
		FROM tables
		WHERE tableid = ?
	`, tableId)

	err := row.Scan(
		&t.ID, &t.Number_of_guests, &t.Table_number, &t.Created_at,
		&t.Updated_at, &t.Table_id,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
func UpdateTable(table *models.Table) error {
	stmt, err := db.DB.Prepare(`
        UPDATE tables SET numberofguests = ?, tablenumber = ?, updatedat = ?, tableid = ?
        WHERE id = ?
    `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(table.Number_of_guests, table.Table_number, table.Updated_at, table.Table_id, table.ID)
	return err
}
func TableExists(tableID string) (bool, error) {
	var id int64
	err := db.DB.QueryRow(`SELECT id FROM tables WHERE tableid = ?`, tableID).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil 
		}
		return false, err 
	}
	return true, nil 
}
