package models

import (
	"time"
)

type Invoice struct {
	ID               int64
	Invoice_id       string
	Order_id         string
	Payment_method   *string
	Payment_status   string
	Payment_due_date time.Time
	Created_at       time.Time
	Updated_at       time.Time
}

type OrderItemView struct {
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

type InvoiceView struct {
	Invoice_id       string          `json:"invoice_id"`
	Order_id         string          `json:"order_id"`
	Payment_method   string          `json:"payment_method"`
	Payment_status   string          `json:"payment_status"`
	Payment_due_date time.Time       `json:"payment_due_date"`
	Payment_due      float64         `json:"payment_due"`
	Table_number     string          `json:"table_number"`
	Order_details    []OrderItemView `json:"order_details"`
}
