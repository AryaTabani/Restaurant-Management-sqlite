package services

import (
	"fmt"
	"time"

	"example.com/m/v2/models"
	"example.com/m/v2/repository"
	"example.com/m/v2/validation"
	"github.com/google/uuid"
)

func CreateOrder(order *models.Order) (*models.Order, error) {
	if err := validation.Validator.Struct(order); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}
	if order.Table_id != nil && *order.Table_id != "" {
		tableExists, err := repository.TableExists(*order.Table_id)
		if err != nil {
			return nil, fmt.Errorf("error checking for table: %w", err)
		}
		if !tableExists {
			return nil, ErrTableNotFound
		}
	}
	order.Created_at = time.Now()
	order.Updated_at = time.Now()
	order.Order_id = uuid.New().String()

	lastID, err := repository.CreateOrder(order)
	if err != nil {
		return nil, fmt.Errorf("could not create order: %w", err)
	}

	order.ID = lastID

	return order, nil
}
