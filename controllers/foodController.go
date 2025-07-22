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

func GetAllFoodsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		foods, err := repository.GetAllFoods()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve foods"})
			return
		}
		c.JSON(http.StatusOK, foods)
	}
}

func GetFoodByIDHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		foodId := c.Param("food_id")

		food, err := repository.GetFoodByID(foodId)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "food not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			}
			return
		}

		c.JSON(http.StatusOK, food)
	}
}

func CreateFoodHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var food models.Food
		if err := c.BindJSON(&food); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
			return
		}

		createdFood, err := services.CreateFood(&food)
		if err != nil {
			switch {
			case errors.Is(err, services.ErrMenuNotFound):
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			case strings.Contains(err.Error(), "validation failed"):
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "an unexpected error occurred"})
			}
			return
		}

		c.JSON(http.StatusCreated, createdFood)
	}
}
func UpdateFoodHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		foodId := c.Param("food_id")

		var foodUpdates models.Food
		if err := c.BindJSON(&foodUpdates); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
			return
		}

		updatedFood, err := services.UpdateFood(foodId, &foodUpdates)
		if err != nil {
			switch {
			case errors.Is(err, services.ErrFoodNotFound):
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			case errors.Is(err, services.ErrMenuNotFound):
				c.JSON(http.StatusBadRequest, gin.H{"error": "the specified menu does not exist"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "an unexpected error occurred"})
			}
			return
		}

		c.JSON(http.StatusOK, updatedFood)
	}
}
