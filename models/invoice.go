package models

import(
	"time"
)

type Invoice struct{
	ID int64
	Invoice_id string
	Order_id string
	Payment_method string
	Payment_status string   `binding:"required"`
	Payment_due_date time.Time
	Created_at time.Time
	Update_at time.Time
}