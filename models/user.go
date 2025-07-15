package models

import (
	"errors"
	"time"

	db "example.com/m/v2/DB"
	"example.com/m/v2/utils"
)

type User struct {
	ID            int64
	First_name    string
	Last_name     string
	Password      string `binding:"required"`
	Email         string `binding:"required"`
	Avatar        string
	Phone         string
	Token         string
	Refresh_Token string
	Created_at    time.Time
	Updated_at    time.Time
	User_id       string
}

func (u User) Save() error {
	query := "INSERT INTO users(firstname,lastname,password,email,avatar,phone,createdat,updatedat,userid) VALUES (?,?,?,?,?,?,?,?,?)"
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()
	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}
	result, err := stmt.Exec(u.First_name, u.Last_name, hashedPassword, u.Email, u.Avatar, u.Phone, u.Created_at, u.Updated_at, u.User_id)
	if err != nil {
		return err
	}
	userId, err := result.LastInsertId()
	u.ID = userId
	return err

}

func (u *User) ValidateCredentials() error {
	query := "SELECT id, password FROM users WHERE email = ?"
	row := db.DB.QueryRow(query, u.Email)
	var retrievedPassword string
	err := row.Scan(&u.ID, &retrievedPassword)
	if err != nil {
		return errors.New("credentials invalid")
	}
	passwordIsvalid := utils.CheckPasswordHash(u.Password, retrievedPassword)
	if !passwordIsvalid {
		return errors.New("credentials invalid")
	}
	return err
}
