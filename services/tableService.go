package services

import (
	"database/sql"
	"fmt"
	"time"

	db "example.com/m/v2/DB"
	"example.com/m/v2/models"
	"example.com/m/v2/repository"
	"example.com/m/v2/validation"
	"github.com/google/uuid"
)

func CreateTable(table *models.Table) (*models.Table, error) {
	if err := validation.Validator.Struct(table); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}
	table.Created_at = time.Now()
	table.Updated_at = time.Now()
	table.Table_id = uuid.New().String()

	lastID, err := repository.CreateTable(table)
	if err != nil {
		return nil, fmt.Errorf("could not create food item: %w", err)
	}

	table.ID = lastID

	return table, nil
}
func IsTableTaken(field, value string) (bool, error) {
	if value == "" {
		return false, nil
	}

	query := "SELECT id FROM tables WHERE " + field + " = ?"
	var id int64
	err := db.DB.QueryRow(query, value).Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
