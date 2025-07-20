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

func GetAllFoods() ([]models.Food, error) {
	query := `
		SELECT id, name, price, Food_image, createdat, updatedat, foodid,menuid 
		FROM foods
	`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var foods []models.Food
	for rows.Next() {
		var f models.Food
		err := rows.Scan(
			&f.ID,
			&f.Name,
			&f.Price,
			&f.Food_image,
			&f.Created_at,
			&f.Update_at,
			&f.Food_id,
			&f.Menu_id,
		)
		if err != nil {
			return nil, err
		}
		foods = append(foods, f)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return foods, nil
}
func GetFoodByID(foodID string) (*models.Food, error) {
	query := `
		SELECT id, name, price, foodimage,createdat, updatedat, foodid,menuid
		FROM foods
		WHERE foodid = ?
	`

	var f models.Food
	err := db.DB.QueryRow(query, foodID).Scan(
		&f.ID,
		&f.Name,
		&f.Price,
		&f.Food_image,
		&f.Created_at,
		&f.Update_at,
		&f.Food_id,
		&f.Menu_id,
	)

	if err != nil {
		return nil, err
	}

	return &f, nil
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
