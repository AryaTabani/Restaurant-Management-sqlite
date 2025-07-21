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
	ErrInvalidTimeSpan   = errors.New("the current time is not within the provided start and end dates")
	ErrMissingTimeStamps = errors.New("both start_date and end_date are required for an update")
)

func inTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}

func CreateMenu(menu *models.Menu) (*models.Menu, error) {
	if err := validation.Validator.Struct(menu); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}
	menu.Created_at = time.Now()
	menu.Updated_at = time.Now()
	menu.Menu_id = uuid.New().String()

	lastID, err := repository.CreateMenu(menu)
	if err != nil {
		return nil, fmt.Errorf("could not create menu: %w", err)
	}
	menu.ID = lastID

	return menu, nil
}
func UpdateMenu(menuID string, updates *models.Menu) (*models.Menu, error) {
	existingMenu, err := repository.GetMenu(menuID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrMenuNotFound
		}
		return nil, err 
	}

	if updates.Start_Date.IsZero() || updates.End_Date.IsZero() {
		return nil, ErrMissingTimeStamps
	}

	if !inTimeSpan(updates.Start_Date, updates.End_Date, time.Now()) {
		return nil, ErrInvalidTimeSpan
	}
	existingMenu.Start_Date = updates.Start_Date
	existingMenu.End_Date = updates.End_Date

	if updates.Name != "" {
		existingMenu.Name = updates.Name
	}
	if updates.Category != "" {
		existingMenu.Category = updates.Category
	}

	existingMenu.Updated_at = time.Now()

	if err := repository.UpdateMenu(existingMenu); err != nil {
		return nil, fmt.Errorf("could not update menu: %w", err)
	}

	return existingMenu, nil
}
