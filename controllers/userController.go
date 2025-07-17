package controllers

import (
	"errors"
	"log"
	"net/http"
	"strings"

	db "example.com/m/v2/DB"
	"example.com/m/v2/models"
	"example.com/m/v2/services"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetAllUsers() ([]models.User, error) {
	query := `
		SELECT id, firstname, lastname, password, email, avatar, phone, createdat, updatedat, userid 
		FROM users
	`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(
			&u.ID,
			&u.First_name,
			&u.Last_name,
			&u.Password,
			&u.Email,
			&u.Avatar,
			&u.Phone,
			&u.Created_at,
			&u.Updated_at,
			&u.User_id,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
func GetUserByID(userID string) (*models.User, error) {
	query := `
		SELECT id, firstname, lastname, password, email, avatar, phone, createdat, updatedat, userid
		FROM users
		WHERE userid = ?
	`

	var u models.User
	err := db.DB.QueryRow(query, userID).Scan(
		&u.ID,
		&u.First_name,
		&u.Last_name,
		&u.Password,
		&u.Email,
		&u.Avatar,
		&u.Phone,
		&u.Created_at,
		&u.Updated_at,
		&u.User_id,
	)

	if err != nil {
		return nil, err
	}

	return &u, nil
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
