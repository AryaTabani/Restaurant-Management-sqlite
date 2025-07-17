// In services/authService.go
package services

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"example.com/m/v2/helpers"
	"example.com/m/v2/models"
	"example.com/m/v2/repository"
	"example.com/m/v2/utils"
	"example.com/m/v2/validation"
)

var (
	ErrEmailExists        = errors.New("this email already exists")
	ErrPhoneExists        = errors.New("this phone number already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

func SignUpUser(user *models.User) (*models.User, error) {
	if err := validation.Validator.Struct(user); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	if emailTaken, err := repository.IsFieldTaken("email", user.Email); err != nil || emailTaken {
		if err != nil {
			return nil, err
		}
		return nil, ErrEmailExists
	}
	if phoneTaken, err := repository.IsFieldTaken("phone", user.Phone); err != nil || phoneTaken {
		if err != nil {
			return nil, err
		}
		return nil, ErrPhoneExists
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("could not hash password: %w", err)
	}
	user.Password = hashedPassword

	now := time.Now()
	user.Created_at = now
	user.Updated_at = now
	user.User_id = fmt.Sprintf("%d-%d", now.Unix(), rand.Intn(1000000))
	token, refreshToken, _ := helpers.GenerateAllTokens(user.Email, user.First_name, user.Last_name, user.User_id)
	user.Token = token
	user.Refresh_Token = refreshToken

	lastID, err := repository.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("could not create user: %w", err)
	}
	user.ID = lastID

	return user, nil
}
func LoginUser(email, password string) (*models.User, error) {
	foundUser, err := repository.GetUserByEmail(email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	passwordIsValid := utils.CheckPasswordHash(password, foundUser.Password)
	if !passwordIsValid {
		return nil, ErrInvalidCredentials
	}

	token, refreshToken, _ := helpers.GenerateAllTokens(
		foundUser.Email, foundUser.First_name, foundUser.Last_name, foundUser.User_id,
	)
	foundUser.Token = token
	foundUser.Refresh_Token = refreshToken

	if err := repository.UpdateUserTokens(foundUser.ID, token, refreshToken); err != nil {
		return nil, fmt.Errorf("could not update session: %w", err)
	}

	return foundUser, nil
}
