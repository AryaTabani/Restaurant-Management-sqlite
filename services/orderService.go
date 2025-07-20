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

var ErrOrderNotFound = errors.New("the specified order was not found")

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

func UpdateOrder(orderID string, updates *models.Order) (*models.Order, error) {
	existingOrder, err := repository.GetOrder(orderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrOrderNotFound
		}
		return nil, err 
	}

	if updates.Table_id != nil {
		tableExists, err := repository.TableExists(*updates.Table_id)
		if err != nil {
			return nil, fmt.Errorf("error checking for table: %w", err)
		}
		if !tableExists {
			return nil, ErrTableNotFound
		}
		existingOrder.Table_id = updates.Table_id
	}

	existingOrder.Updated_at = time.Now()

	if err := repository.UpdateOrder(existingOrder); err != nil {
		return nil, fmt.Errorf("could not save updated order: %w", err)
	}

	return existingOrder, nil
}
