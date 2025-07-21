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

func GetAllMenus() ([]models.Menu, error) {
	query := `
		SELECT id, name, category,startdate, enddate, createdat, updatedat, menuid 
		FROM orders
	`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var menus []models.Menu
	for rows.Next() {
		var m models.Menu
		err := rows.Scan(
			&m.ID,
			&m.Name,
			&m.Category,
			&m.Start_Date,
			&m.End_Date,
			&m.Created_at,
			&m.Updated_at,
			&m.Menu_id,
		)
		if err != nil {
			return nil, err
		}
		menus = append(menus, m)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return menus, nil
}
func GetMenuByID(menuID string) (*models.Menu, error) {
	query := `
		SELECT id, name,category, startdate, enddate, createdat, updatedat, menuid
		FROM menus
		WHERE menuid = ?
	`

	var m models.Menu
	err := db.DB.QueryRow(query, menuID).Scan(
		&m.ID,
		&m.Name,
		&m.Category,
		&m.Start_Date,
		&m.End_Date,
		&m.Created_at,
		&m.Updated_at,
		&m.Menu_id,
	)

	if err != nil {
		return nil, err
	}

	return &m, nil
}

func CreateMenuHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var menu models.Menu
		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
			return
		}

		createdMenu, err := services.CreateMenu(&menu)
		if err != nil {
			switch {
			case strings.Contains(err.Error(), "validation failed"):
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "an unexpected error occurred"})
			}
			return
		}

		c.JSON(http.StatusCreated, createdMenu)
	}
}

func UpdateMenuHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		menuId := c.Param("menu_id")

		var menuUpdates models.Menu
		if err := c.BindJSON(&menuUpdates); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
			return
		}

		updatedMenu, err := services.UpdateMenu(menuId, &menuUpdates)
		if err != nil {
			switch {
			case errors.Is(err, services.ErrMenuNotFound):
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			case errors.Is(err, services.ErrInvalidTimeSpan), errors.Is(err, services.ErrMissingTimeStamps):
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "an unexpected error occurred"})
			}
			return
		}

		c.JSON(http.StatusOK, updatedMenu)
	}
}
