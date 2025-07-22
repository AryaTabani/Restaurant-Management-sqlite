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

func GetAllTablesHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		tables, err := repository.GetAllTables()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve tables"})
			return
		}
		c.JSON(http.StatusOK, tables)
	}
}

func GetTableByIDHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		tableId := c.Param("table_id")

		table, err := repository.GetTableByID(tableId)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "table not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			}
			return
		}

		c.JSON(http.StatusOK, table)
	}
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
func UpdateTableHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		tableID := c.Param("Table_id")

		var tableUpdates models.Table
		if err := c.BindJSON(&tableUpdates); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
			return
		}

		updatedTable, err := services.UpdateTable(tableID, &tableUpdates)
		if err != nil {
			switch {
			case errors.Is(err, services.ErrTableNotFound):
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "an unexpected error occurred"})
			}
			return
		}

		c.JSON(http.StatusOK, updatedTable)
	}
}
