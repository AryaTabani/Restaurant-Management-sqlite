package services

import (
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

var (
	ErrTableNotFound = errors.New("the specified table was not found")
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
		return nil, fmt.Errorf("could not create table: %w", err)
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
func UpdateTable(tableID string, updates *models.Table) (*models.Table, error) {
	existingTable, err := repository.GetTableById(tableID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrTableNotFound
		}
		return nil, err
	}

	if updates.Number_of_guests != nil {
		existingTable.Number_of_guests = updates.Number_of_guests
	}
	if updates.Table_number != nil {
		existingTable.Table_number = updates.Table_number
	}
	existingTable.Updated_at = time.Now()

	if err := repository.UpdateTable(existingTable); err != nil {
		return nil, fmt.Errorf("could not save updated table: %w", err)
	}

	return existingTable, nil
}
