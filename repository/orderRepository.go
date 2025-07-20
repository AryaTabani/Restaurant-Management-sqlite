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
func GetOrder(orderID string) (*models.Order, error) {
	var order models.Order
	err := db.DB.QueryRow(`
		SELECT id, orderdate, createdat, updatedat, orderid, tableid
		FROM orders WHERE orderid = ?`, orderID).Scan(
		&order.ID, &order.Order_date, &order.Created_at, &order.Updated_at,
		&order.Order_id, &order.Table_id,
	)
	if err != nil {
		return nil, err 
	}
	return &order, nil
}


func UpdateOrder(order *models.Order) error {
	stmt, err := db.DB.Prepare(`
		UPDATE orders SET table_id = ?, updated_at = ?
		WHERE id = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(order.Table_id, order.Updated_at, order.ID)
	return err
}