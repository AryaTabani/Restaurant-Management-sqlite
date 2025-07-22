package controllers

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"example.com/m/v2/models"
	"example.com/m/v2/repository"
	"example.com/m/v2/services"
	"github.com/gin-gonic/gin"
)

func GetAllOrdersHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		orders, err := repository.GetAllOrders()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve orders"})
			return
		}
		c.JSON(http.StatusOK, orders)
	}
}

func GetOrderByIDHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderId := c.Param("order_id")

		order, err := repository.GetOrderByID(orderId)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			}
			return
		}

		c.JSON(http.StatusOK, order)
	}
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
