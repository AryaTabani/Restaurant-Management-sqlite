package controllers

import (
	"net/http"
	"strings"

	db "example.com/m/v2/DB"
	"example.com/m/v2/models"
	"example.com/m/v2/services"
	"github.com/gin-gonic/gin"
)

func GetAllTables() ([]models.Table, error) {
	query := `
		SELECT id, numberofguests, tablenumber, createdat, updatedat, tableid 
		FROM tables
	`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []models.Table
	for rows.Next() {
		var t models.Table
		err := rows.Scan(
			&t.ID,
			&t.Number_of_guests,
			&t.Table_number,
			&t.Created_at,
			&t.Updated_at,
			&t.Table_id,
		)
		if err != nil {
			return nil, err
		}
		tables = append(tables, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}

func GetTableByID(tableID string) (*models.Table, error) {
	query := `
		SELECT id, numberofguests, tablenumber, createdat, updatedat, tableid 
		FROM tables
		WHERE tableid = ?
	`

	var t models.Table
	err := db.DB.QueryRow(query, tableID).Scan(
		&t.ID,
		&t.Number_of_guests,
		&t.Table_number,
		&t.Created_at,
		&t.Updated_at,
		&t.Table_id,
	)

	if err != nil {
		return nil, err
	}

	return &t, nil
}
func CreateTableHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var table models.Table
		if err := c.BindJSON(&table); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
			return
		}

		createdTable, err := services.CreateTable(&table)
		if err != nil {
			switch {
			// case errors.Is(err, services.IsTableTaken("tableid", table.Table_id)):
			// 	c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			case strings.Contains(err.Error(), "validation failed"):
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "an unexpected error occurred"})
			}
			return
		}

		c.JSON(http.StatusCreated, createdTable)
	}
}
