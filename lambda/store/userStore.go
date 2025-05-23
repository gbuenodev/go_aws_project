package store

import (
	"fmt"
	"lambda_func/utils"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	TABLE_NAME = "userTable"
)

type RegisterUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password"`
}

type UserStore interface {
	DoesUserExist(username string) (bool, error)
	InsertUser(user User) error
	GetUser(username string) (User, error)
}

func NewUser(registerUser RegisterUser) (User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerUser.Password), 10)
	if err != nil {
		return User{}, err
	}

	return User{
		Username:     registerUser.Username,
		PasswordHash: string(hashedPassword),
	}, nil
}

func ValidatePassword(hashedPassword, plaintextPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plaintextPassword))
	return err == nil
}

func CreateToken(user User) (string, error) {
	secretArn := os.Getenv("JWT_SECRET_ARN")
	if secretArn == "" {
		return "", fmt.Errorf("JWT_SECRET_ARN is not set")
	}

	jwtSecret, err := utils.FetchJWTSecret(secretArn)
	if err != nil {
		return "", fmt.Errorf("failed to fetch JWT secret: %w", err)
	}

	now := time.Now()
	validUntil := now.Add(time.Hour * 1).Unix()

	claims := jwt.MapClaims{
		"user":    user.Username,
		"expires": validUntil,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
