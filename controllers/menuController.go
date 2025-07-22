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

func GetAllMenusHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		menus, err := repository.GetAllMenus()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve menus"})
			return
		}
		c.JSON(http.StatusOK, menus)
	}
}

func GetMenuByIDHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		menuId := c.Param("menu_id")

		menu, err := repository.GetMenuByID(menuId)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "menu not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			}
			return
		}

		c.JSON(http.StatusOK, menu)
	}
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
