package controllers

import (
	"errors"
	"net/http"
	"strings"

	"example.com/m/v2/models"
	"example.com/m/v2/services"
	"github.com/gin-gonic/gin"
)

func CreateOrderWithItemsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var orderItemPack models.OrderItemPack
		if err := c.BindJSON(&orderItemPack); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		createdOrder, err := services.CreateOrderWithItems(orderItemPack)
		if err != nil {
			if strings.Contains(err.Error(), "validation failed") {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create order"})
			}
			return
		}
		c.JSON(http.StatusCreated, createdOrder)
	}
}

func GetAllOrderItemsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		items, err := services.GetAllOrderItems()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve order items"})
			return
		}
		c.JSON(http.StatusOK, items)
	}
}

func GetOrderItemHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderItemId := c.Param("order_item_id")

		item, err := services.GetOrderItem(orderItemId)
		if err != nil {
			if errors.Is(err, services.ErrOrderItemNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			}
			return
		}
		c.JSON(http.StatusOK, item)
	}
}

func GetOrderItemsByOrderHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderId := c.Param("order_id")

		items, err := services.GetOrderItemsByOrderID(orderId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve order items for the specified order"})
			return
		}
		c.JSON(http.StatusOK, items)
	}
}

func UpdateOrderItemHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderItemId := c.Param("order_item_id")
		var updates models.OrderItem

		if err := c.BindJSON(&updates); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		updatedItem, err := services.UpdateOrderItem(orderItemId, &updates)
		if err != nil {
			if errors.Is(err, services.ErrOrderItemNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update order item"})
			}
			return
		}
		c.JSON(http.StatusOK, updatedItem)
	}
}
