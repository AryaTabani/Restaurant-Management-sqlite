package models

import (
	"time"
)

type Order struct {
	ID         int64
	Order_date time.Time `binding:"required"`
	Created_at time.Time
	Updated_at time.Time
	Order_id   string
	Table_id   *string `binding:"required"`
}
