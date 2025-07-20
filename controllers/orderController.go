package controllers

import (
	"errors"
	"net/http"
	"strings"

	db "example.com/m/v2/DB"
	"example.com/m/v2/models"
	"example.com/m/v2/services"
	"github.com/gin-gonic/gin"
)

func GetAllOrders() ([]models.Order, error) {
	query := `
		SELECT id, orderdate, createdat, updatedat, orderid,tableid 
		FROM orders
	`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		err := rows.Scan(
			&o.ID,
			&o.Order_date,
			&o.Created_at,
			&o.Updated_at,
			&o.Order_id,
			&o.Table_id,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
func GetOrderByID(orderID string) (*models.Order, error) {
	query := `
		SELECT id, orderdate,createdat, updatedat, orderid,tableid
		FROM orders
		WHERE orderid = ?
	`

	var o models.Order
	err := db.DB.QueryRow(query, orderID).Scan(
		&o.ID,
		&o.Order_date,
		&o.Created_at,
		&o.Updated_at,
		&o.Order_id,
		&o.Table_id,
	)

	if err != nil {
		return nil, err
	}

	return &o, nil
}
func CreateOrderHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var order models.Order
		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
			return
		}

		createdOrder, err := services.CreateOrder(&order)
		if err != nil {
			switch {
			case errors.Is(err, services.ErrTableNotFound):
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			case strings.Contains(err.Error(), "validation failed"):
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "an unexpected error occurred"})
			}
			return
		}

		c.JSON(http.StatusCreated, createdOrder)
	}
}

func UpdateOrderHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderID := c.Param("Order_id")

		var orderUpdates models.Order
		if err := c.BindJSON(&orderUpdates); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
			return
		}

		updatedOrder, err := services.UpdateOrder(orderID, &orderUpdates)
		if err != nil {
			switch {
			case errors.Is(err, services.ErrOrderNotFound):
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			case errors.Is(err, services.ErrTableNotFound):
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "an unexpected error occurred"})
			}
			return
		}

		c.JSON(http.StatusOK, updatedOrder)
	}
}
