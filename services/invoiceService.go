// In services/invoiceService.go
package services

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"example.com/m/v2/models"
	"example.com/m/v2/repository"
	"example.com/m/v2/validation"
	"github.com/google/uuid"
)

var ErrInvoiceNotFound = errors.New("the specified invoice was not found")

func GetAllInvoices() ([]models.Invoice, error) {
	return repository.GetAllInvoices()
}

func GetInvoice(invoiceID string) (*models.InvoiceView, error) {
	invoice, err := repository.GetInvoiceDetails(invoiceID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrInvoiceNotFound
		}
		return nil, err
	}
	return invoice, nil
}


func CreateInvoice(invoice *models.Invoice) (*models.Invoice, error) {
	orderExists, err := repository.OrderExists(invoice.Order_id)
	if err != nil {
		return nil, fmt.Errorf("error checking for order: %w", err)
	}
	if !orderExists {
		return nil, ErrOrderNotFound
	}

	if invoice.Payment_status == "" {
		invoice.Payment_status = "PENDING"
	}
	now := time.Now()
	invoice.Payment_due_date = now.AddDate(0, 0, 1) 
	invoice.Created_at = now
	invoice.Updated_at = now
	invoice.Invoice_id = uuid.New().String()

	if err := validation.Validator.Struct(invoice); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	lastID, err := repository.CreateInvoice(invoice)
	if err != nil {
		return nil, fmt.Errorf("could not create invoice: %w", err)
	}
	invoice.ID = lastID

	return invoice, nil
}

func UpdateInvoice(invoiceID string, updates *models.Invoice) (*models.Invoice, error) {
	existingInvoice, err := repository.GetInvoice(invoiceID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrInvoiceNotFound
		}
		return nil, err
	}

	if updates.Payment_method != nil {
		existingInvoice.Payment_method = updates.Payment_method
	}
	if updates.Payment_status != "" {
		existingInvoice.Payment_status = updates.Payment_status
	}

	existingInvoice.Updated_at = time.Now()

	if err := repository.UpdateInvoice(existingInvoice); err != nil {
		return nil, fmt.Errorf("could not update invoice: %w", err)
	}
	return existingInvoice, nil
}
