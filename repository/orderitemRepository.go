package repository

import (
	"database/sql"

	db "example.com/m/v2/DB"
	"example.com/m/v2/models"
)

func CreateOrderItemInTx(tx *sql.Tx, item models.OrderItem) error {
	stmt, err := tx.Prepare(`
		INSERT INTO orderitems (quantity, unitprice, createdat, updatedat, foodid, orderitemid, orderid)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(item.Quantity, item.Unit_price, item.Created_at, item.Updated_at, item.Food_id, item.Order_item_id, item.Order_id)
	return err
}

func GetOrderItem(orderItemID string) (*models.OrderItem, error) {
	var item models.OrderItem
	err := db.DB.QueryRow(`SELECT * FROM orderitems WHERE orderitemid = ?`, orderItemID).Scan(
		&item.ID, &item.Quantity, &item.Unit_price, &item.Created_at, &item.Updated_at,
		&item.Food_id, &item.Order_item_id, &item.Order_id,
	)
	return &item, err
}

func GetAllOrderItems() ([]models.OrderItem, error) {
	rows, err := db.DB.Query("SELECT * FROM orderitems")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		err := rows.Scan(
			&item.ID, &item.Quantity, &item.Unit_price, &item.Created_at, &item.Updated_at,
			&item.Food_id, &item.Order_item_id, &item.Order_id,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}
func GetOrderItemsByOrderID(orderID string) ([]models.OrderItemView, error) {
	query := `
		SELECT f.name, oi.quantity, f.price
		FROM orderitems oi
		LEFT JOIN foods f ON oi.foodid = f.foodid
		WHERE oi.orderid = ?`

	rows, err := db.DB.Query(query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.OrderItemView
	for rows.Next() {
		var item models.OrderItemView
		if err := rows.Scan(&item.Name, &item.Quantity, &item.Price); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}
func UpdateOrderItem(item *models.OrderItem) error {
	stmt, err := db.DB.Prepare(`
		UPDATE orderitems SET quantity = ?, unitprice = ?, updatedat = ?
		WHERE id = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(item.Quantity, item.Unit_price, item.Updated_at, item.ID)
	return err
}