package repository

import (
	db "example.com/m/v2/DB"
	"example.com/m/v2/models"
)

func CreateOrder(order *models.Order) (int64, error) {
	stmt, err := db.DB.Prepare(`
		INSERT INTO orders (orderdate, createdat, updatedat, orderid, tableid)
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		order.Order_date,
		order.Created_at,
		order.Updated_at,
		order.Order_id,
		order.Table_id,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}
