package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"time"

	db "example.com/m/v2/DB"
	"example.com/m/v2/models"
	"example.com/m/v2/repository"
	"example.com/m/v2/validation"
	"github.com/google/uuid"
)

var ErrOrderItemNotFound = errors.New("the specified order item was not found")

func CreateOrderWithItems(orderPack models.OrderItemPack) (*models.Order, error) {
	tx, err := db.DB.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback()

	var order models.Order
	now := time.Now()
	order.Order_date = now
	order.Created_at = now
	order.Updated_at = now
	order.Order_id = uuid.New().String()
	order.Table_id = orderPack.Table_id

	orderID, err := repository.CreateOrderInTx(tx, &order)
	if err != nil {
		return nil, fmt.Errorf("could not create order record: %w", err)
	}
	order.ID = orderID

	for _, item := range orderPack.Order_items {
		item.Order_id = order.Order_id
		item.Created_at = now
		item.Updated_at = now
		item.Order_item_id = uuid.New().String()

		if err := validation.Validator.Struct(item); err != nil {
			return nil, fmt.Errorf("validation failed for order item: %w", err)
		}
		if err := repository.CreateOrderItemInTx(tx, item); err != nil {
			return nil, fmt.Errorf("could not create order item record: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("could not commit transaction: %w", err)
	}
	return &order, nil
}

func GetAllOrderItems() ([]models.OrderItem, error) {
	return repository.GetAllOrderItems()
}

func GetOrderItem(orderItemID string) (*models.OrderItem, error) {
	item, err := repository.GetOrderItem(orderItemID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrOrderItemNotFound
		}
		return nil, err
	}
	return item, nil
}

func GetOrderItemsByOrderID(orderID string) ([]models.OrderItemView, error) {
	return repository.GetOrderItemsByOrderID(orderID)
}

func UpdateOrderItem(orderItemID string, updates *models.OrderItem) (*models.OrderItem, error) {
	existingItem, err := repository.GetOrderItem(orderItemID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrOrderItemNotFound
		}
		return nil, err
	}

	if updates.Quantity > 0 {
		existingItem.Quantity = updates.Quantity
	}
	if updates.Unit_price > 0 {
		existingItem.Unit_price = updates.Unit_price
	}

	existingItem.Updated_at = time.Now()

	if err := repository.UpdateOrderItem(existingItem); err != nil {
		return nil, fmt.Errorf("could not update order item: %w", err)
	}
	return existingItem, nil
}
