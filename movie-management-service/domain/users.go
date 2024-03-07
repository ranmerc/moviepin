package domain

import (
	"database/sql"
	"errors"
	"movie-management-service/utils"

	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrInvalidCredentials is the error message when the credentials are invalid.
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrFailedLogin is the error message when the login fails.
	ErrFailedLogin = errors.New("failed to login")

	// ErrUsernameExists is the error message when the username already exists.
	ErrUsernameExists = errors.New("username already exists")

	// ErrFailedRegister is the error message when the user register fails.
	ErrFailedRegister = errors.New("failed to sign up")
)

func (ms MovieService) LoginUser(username, password string) error {
	sqlStatement := `
		SELECT "password"
		FROM "users"
		WHERE "username" = $1;
	`

	row := ms.db.QueryRow(sqlStatement, username)

	var hashedPassword string
	if err := row.Scan(&hashedPassword); err != nil {
		if err == sql.ErrNoRows {
			return ErrInvalidCredentials
		}

		utils.ErrorLogger.Println(err)
		return ErrFailedLogin
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return ErrInvalidCredentials
	}

	return nil
}

func (ms MovieService) RegisterUser(username, password string) error {
	sqlStatement := `
		SELECT COUNT(*) 
		FROM "users" 
		WHERE "username" = $1;
	`

	row := ms.db.QueryRow(sqlStatement, username)

	var count int

	if err := row.Scan(&count); err != nil && err != sql.ErrNoRows {
		utils.ErrorLogger.Println(err)
		return ErrFailedRegister
	}

	if count > 0 {
		return ErrUsernameExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		utils.ErrorLogger.Println(err)
		return ErrFailedRegister
	}

	sqlStatement = `
		INSERT INTO "users" ("username", "password")
		VALUES ($1, $2);
	`

	if _, err := ms.db.Exec(sqlStatement, username, hashedPassword); err != nil {
		utils.ErrorLogger.Println(err)
		return ErrFailedRegister
	}

	return nil
}
