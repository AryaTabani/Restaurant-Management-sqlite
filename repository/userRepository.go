// In repository/userRepository.go
package repository

import (
	"database/sql"

	db "example.com/m/v2/DB"
	"example.com/m/v2/models"
)

func IsFieldTaken(field, value string) (bool, error) {
	if value == "" {
		return false, nil
	}

	query := "SELECT id FROM users WHERE " + field + " = ?"
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

func CreateUser(user *models.User) (int64, error) {
	stmt, err := db.DB.Prepare(`
		INSERT INTO users (firstname, lastname, password, email, avatar, phone, createdat, updatedat, userid)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		user.First_name, user.Last_name, user.Password, user.Email,
		user.Avatar, user.Phone, user.Created_at, user.Updated_at, user.User_id,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func GetUserByEmail(email string) (*models.User, error) {
	var u models.User
	row := db.DB.QueryRow(`
		SELECT id, firstname, lastname, password, email, avatar, phone, createdat, updatedat, userid
		FROM users
		WHERE email = ?
	`, email)

	err := row.Scan(
		&u.ID, &u.First_name, &u.Last_name, &u.Password, &u.Email,
		&u.Avatar, &u.Phone, &u.Created_at, &u.Updated_at, &u.User_id,
	)
	if err != nil {
		return nil, err 
	}
	return &u, nil
}

func UpdateUserTokens(id int64, token, refreshToken string) error {
	_, err := db.DB.Exec(`UPDATE users SET token = ?, refreshtoken = ? WHERE id = ?`, token, refreshToken, id)
	return err
}

