package models

import (
	"time"
)

type OrderItem struct {
	ID            int64
	Quantity      int     `json:"quantity"`
	Unit_price    float64 `json:"unit_price"`
	Created_at    time.Time
	Updated_at    time.Time
	Food_id       string `json:"food_id"`
	Order_item_id string
	Order_id      string `json:"order_id"`
}

type OrderItemPack struct {
	Table_id    *string     `json:"table_id"`
	Order_items []OrderItem `json:"order_items"`
}
