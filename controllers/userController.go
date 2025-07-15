package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	db "example.com/m/v2/DB"
	"example.com/m/v2/models"
	"example.com/m/v2/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}

		offset := (page - 1) * recordPerPage

		var totalCount int
		err = db.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&totalCount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not count users"})
			return
		}

		rows, err := db.DB.Query(`
            SELECT id, firstname, lastname, password, email, avatar, phone, createdat, updatedat, userid
            FROM users
            LIMIT ? OFFSET ?
        `, recordPerPage, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch users"})
			return
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
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error scanning user"})
				return
			}
			users = append(users, u)
		}

		c.JSON(http.StatusOK, gin.H{
			"total_count":   totalCount,
			"page":          page,
			"recordPerPage": recordPerPage,
			"users":         users,
		})
	}
}
func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")

		query := `
            SELECT id, firstname, lastname, password, email, avatar, phone, createdat, updatedat, userid
            FROM users
            WHERE userid = ?
        `

		var u models.User
		err := db.DB.QueryRow(query, userId).Scan(
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
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while fetching user"})
			}
			return
		}

		u.Password = ""

		c.JSON(http.StatusOK, u)
	}
}
func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if validationErr := validate.Struct(user); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		var existingID int64
		err := db.DB.QueryRow(`SELECT id FROM users WHERE email = ?`, user.Email).Scan(&existingID)
		if err != sql.ErrNoRows {
			if err == nil {
				c.JSON(http.StatusConflict, gin.H{"error": "this email already exists"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error checking email"})
			return
		}

		if user.Phone != "" {
			err = db.DB.QueryRow(`SELECT id FROM users WHERE phone = ?`, user.Phone).Scan(&existingID)
			if err != sql.ErrNoRows {
				if err == nil {
					c.JSON(http.StatusConflict, gin.H{"error": "this phone number already exists"})
					return
				}
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error checking phone"})
				return
			}
		}

		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error hashing password"})
			return
		}
		user.Password = hashedPassword

		now := time.Now()
		user.Created_at = now
		user.Updated_at = now
		user.User_id = fmt.Sprintf("%d-%d", now.Unix(), rand.Intn(1000000)) // or use uuid

		token, refreshToken, _ := helpers.GenerateAllTokens(user.Email, user.First_name, user.Last_name, user.User_id)
		user.Token = token
		user.Refresh_Token = refreshToken

		stmt, err := db.DB.Prepare(`
            INSERT INTO users (firstname, lastname, password, email, avatar, phone, createdat, updatedat, userid)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
        `)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error preparing insert"})
			return
		}
		defer stmt.Close()

		res, err := stmt.Exec(
			user.First_name,
			user.Last_name,
			user.Password,
			user.Email,
			user.Avatar,
			user.Phone,
			user.Created_at,
			user.Updated_at,
			user.User_id,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error inserting user"})
			return
		}

		lastID, _ := res.LastInsertId()
		user.ID = lastID

		user.Password = ""

		c.JSON(http.StatusOK, user)
	}
}
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginRequest models.User
		if err := c.BindJSON(&loginRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var foundUser models.User
		row := db.DB.QueryRow(`
            SELECT id, firstname, lastname, password, email, avatar, phone, createdat, updatedat, userid
            FROM users
            WHERE email = ?
        `, loginRequest.Email)

		err := row.Scan(
			&foundUser.ID,
			&foundUser.First_name,
			&foundUser.Last_name,
			&foundUser.Password,
			&foundUser.Email,
			&foundUser.Avatar,
			&foundUser.Phone,
			&foundUser.Created_at,
			&foundUser.Updated_at,
			&foundUser.User_id,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error finding user"})
			}
			return
		}

		passwordIsValid := utils.CheckPasswordHash(loginRequest.Password, foundUser.Password)
		if !passwordIsValid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}

		token, refreshToken, _ := helpers.GenerateAllTokens(
			foundUser.Email,
			foundUser.First_name,
			foundUser.Last_name,
			foundUser.User_id,
		)

		foundUser.Token = token
		foundUser.Refresh_Token = refreshToken

		_, _ = db.DB.Exec(`UPDATE users SET token = ?, refreshtoken = ? WHERE id = ?`, token, refreshToken, foundUser.ID)

		foundUser.Password = ""

		c.JSON(http.StatusOK, foundUser)
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
