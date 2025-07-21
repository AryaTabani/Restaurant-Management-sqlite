package models

import (
	"time"
)

type Menu struct {
	ID         int64
	Name       string `binding:"required"`
	Category   string `binding:"required"`
	Start_Date time.Time
	End_Date   time.Timer
	Created_at time.Time
	Updated_at  time.Time
	Menu_id    string
}
