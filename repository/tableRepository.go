package repository

import (
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
