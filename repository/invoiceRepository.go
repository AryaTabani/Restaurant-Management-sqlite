package repository

import (
	"database/sql"

	db "example.com/m/v2/DB"
	"example.com/m/v2/models"
)

func GetAllInvoices() ([]models.Invoice, error) {
	rows, err := db.DB.Query(`SELECT id, invoiceid, orderid, paymentmethod, paymentstatus, paymentduedate FROM invoices`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []models.Invoice
	for rows.Next() {
		var inv models.Invoice
		if err := rows.Scan(&inv.ID, &inv.Invoice_id, &inv.Order_id, &inv.Payment_method, &inv.Payment_status, &inv.Payment_due_date); err != nil {
			return nil, err
		}
		invoices = append(invoices, inv)
	}
	return invoices, nil
}

func GetInvoiceDetails(invoiceID string) (*models.InvoiceView, error) {
	query := `
		SELECT
			i.invoiceid, i.orderid, i.paymentmethod, i.paymentstatus, i.paymentduedate,
			o.paymentdue, t.tablenumber,
			f.name, oi.quantity, f.price
		FROM invoices i
		LEFT JOIN orders o ON i.orderid = o.orderid
		LEFT JOIN tables t ON o.tableid = t.tableid
		LEFT JOIN order_items oi ON o.orderid = oi.orderid
		LEFT JOIN foods f ON oi.foodid = f.foodid
		WHERE i.invoice_id = ?`

	rows, err := db.DB.Query(query, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoiceView models.InvoiceView
	isFirstRow := true

	for rows.Next() {
		var item models.OrderItemView
		var tableNumber, paymentMethod sql.NullString
		var paymentDue sql.NullFloat64

		err := rows.Scan(
			&invoiceView.Invoice_id, &invoiceView.Order_id, &paymentMethod, &invoiceView.Payment_status, &invoiceView.Payment_due_date,
			&paymentDue, &tableNumber,
			&item.Name, &item.Quantity, &item.Price,
		)
		if err != nil {
			return nil, err
		}

		if isFirstRow {
			invoiceView.Table_number = tableNumber.String
			invoiceView.Payment_method = paymentMethod.String
			invoiceView.Payment_due = paymentDue.Float64
			isFirstRow = false
		}
		invoiceView.Order_details = append(invoiceView.Order_details, item)
	}

	if isFirstRow {
		return nil, sql.ErrNoRows
	}

	return &invoiceView, nil
}
func CreateInvoice(invoice *models.Invoice) (int64, error) {
	stmt, err := db.DB.Prepare(`
		INSERT INTO invoices (invoiceid, orderid, paymentmethod, paymentstatus, paymentduedate, createdat, updatedat)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		invoice.Invoice_id, invoice.Order_id, invoice.Payment_method, invoice.Payment_status,
		invoice.Payment_due_date, invoice.Created_at, invoice.Updated_at,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetInvoice(invoiceID string) (*models.Invoice, error) {
	var invoice models.Invoice
	err := db.DB.QueryRow(`
		SELECT id, invoiceid, orderid, paymentmethod, paymentstatus, paymentduedate, createdat, updatedat
		FROM invoices WHERE invoice_id = ?`, invoiceID).Scan(
		&invoice.ID, &invoice.Invoice_id, &invoice.Order_id, &invoice.Payment_method, &invoice.Payment_status,
		&invoice.Payment_due_date, &invoice.Created_at, &invoice.Updated_at,
	)
	if err != nil {
		return nil, err
	}
	return &invoice, nil
}

func UpdateInvoice(invoice *models.Invoice) error {
	stmt, err := db.DB.Prepare(`
		UPDATE invoices SET paymentmethod = ?, paymentstatus = ?, updatedat = ?
		WHERE id = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(invoice.Payment_method, invoice.Payment_status, invoice.Updated_at, invoice.ID)
	return err
}