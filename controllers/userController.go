package controllers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"

	"example.com/m/v2/models"
	"example.com/m/v2/repository"
	"example.com/m/v2/services"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetAllUsersHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := repository.GetAllUsers()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve users"})
			return
		}
		c.JSON(http.StatusOK, users)
	}
}

func GetUserByIDHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")

		user, err := repository.GetUserByID(userId)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			}
			return
		}

		user.Password = ""
		c.JSON(http.StatusOK, user)
	}
}
func SignUpHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
			return
		}

		createdUser, err := services.SignUpUser(&user)
		if err != nil {
			switch {
			case errors.Is(err, services.ErrEmailExists), errors.Is(err, services.ErrPhoneExists):
				c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			case strings.Contains(err.Error(), "validation failed"):
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			default:

				c.JSON(http.StatusInternalServerError, gin.H{"error": "an unexpected error occurred"})
			}
			return
		}

		createdUser.Password = ""
		c.JSON(http.StatusOK, createdUser)
	}
}
func LoginHandler() gin.HandlerFunc {
	type loginRequest struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	return func(c *gin.Context) {
		var req loginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
			return
		}

		loggedInUser, err := services.LoginUser(req.Email, req.Password)
		if err != nil {
			if errors.Is(err, services.ErrInvalidCredentials) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "an unexpected error occurred"})
			}
			return
		}

		loggedInUser.Password = ""
		c.JSON(http.StatusOK, loggedInUser)
	}
}
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(plainPassword, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
