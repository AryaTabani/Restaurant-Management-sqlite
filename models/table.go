package models

import (
	"time"
)

type Table struct {
	ID               int64
	Number_of_guests *int `binding:"required"`
	Table_number     *int `binding:"required"`
	Created_at       time.Time
	Updated_at       time.Time
	Table_id         string
}
