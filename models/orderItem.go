package models

import (
	"time"
)

type OrderItem struct {
	ID            int64
	Quantity      string  `binding:"required"`
	Unit_price    float64 `binding:"required"`
	Created_at    time.Time
	Updated_at    time.Time
	Food_id       string `binding:"required"`
	Order_item_id string
	Order_id      string `binding:"required"`
}
