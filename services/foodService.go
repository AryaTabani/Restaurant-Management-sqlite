package services

import (
	"errors"
	"fmt"
	"time"

	"example.com/m/v2/models"
	"example.com/m/v2/repository"
	"example.com/m/v2/validation"
	"github.com/google/uuid"
)

var (
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
