package models

import (
	"time"
)

type User struct {
	ID            int64
	First_name    string
	Last_name     string
	Password      string `binding:"required"`
	Email         string `binding:"required"`
	Avatar        string
	Phone         string
	Token         string
	Refresh_Token string
	Created_at    time.Time
	Updated_at    time.Time
	User_id       string
}
