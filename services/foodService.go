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

var (
	ErrFoodNotFound = errors.New("the specified food item was not found")
	ErrMenuNotFound = errors.New("the specified menu does not exist")
)

func CreateFood(food *models.Food) (*models.Food, error) {
	if err := validation.Validator.Struct(food); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	menuExists, err := repository.MenuExists(food.Menu_id)
	if err != nil {
		return nil, fmt.Errorf("error checking for menu: %w", err)
	}
	if !menuExists {
		return nil, ErrMenuNotFound
	}

	food.Created_at = time.Now()
	food.Update_at = time.Now()
	food.Food_id = uuid.New().String()

	lastID, err := repository.CreateFood(food)
	if err != nil {
		return nil, fmt.Errorf("could not create food item: %w", err)
	}

	food.ID = lastID

	return food, nil
}
func UpdateFood(foodID string, updates *models.Food) (*models.Food, error) {
	existingFood, err := repository.GetFoodById(foodID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrFoodNotFound
		}
		return nil, err
	}

	if updates.Name != "" {
		existingFood.Name = updates.Name
	}
	if updates.Price != 0 {
		existingFood.Price = updates.Price
	}
	if updates.Food_image != "" {
		existingFood.Food_image = updates.Food_image
	}
	if updates.Menu_id != "" && updates.Menu_id != existingFood.Menu_id {
		menuExists, err := repository.MenuExists(updates.Menu_id)
		if err != nil {
			return nil, fmt.Errorf("error checking for menu: %w", err)
		}
		if !menuExists {
			return nil, ErrMenuNotFound
		}
		existingFood.Menu_id = updates.Menu_id
	}

	existingFood.Update_at = time.Now()

	if err := repository.UpdateFood(existingFood); err != nil {
		return nil, fmt.Errorf("could not save updated food: %w", err)
	}

	return existingFood, nil
}
