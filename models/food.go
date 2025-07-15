package models

import (
	"time"
)

type Food struct {
	ID         int64
	Name       string  `binding:"required"`
	Price      float64 `binding:"required"`
	Food_image string  `binding:"required"`
	Created_at time.Time
	Update_at  time.Time
	Food_id    string
	Menu_id    string `binding:"required"`
}
