package services

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"resons/v0/api/api/db"

	"github.com/labstack/echo/v4"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"time"
)

type Token struct {
	ID        string
	UserID    uuid.UUID
	ExpiresAt time.Time
}

type AuthService struct {
	db *db.Queries
}

func NewAuthService(queries *db.Queries) *AuthService {
	return &AuthService{
		db: queries,
	}
}

func (ac *AuthService) UpsertUserWithToken(c echo.Context, tokenString string, email string) error {
	expiration := time.Now().Add(24 * time.Hour)
	magicToken := sql.NullString{String: tokenString, Valid: true}
	expirationTime := sql.NullTime{Time: expiration, Valid: true}

	if err := ac.db.UpsertUserWithToken(context.Background(), db.UpsertUserWithTokenParams{
		Email:           email,
		MagicToken:      magicToken,
		TokenExpiration: expirationTime,
	}); err != nil {
		log.Printf("Database error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Unable to upsert user"})
	}

	return nil
}

func (ac *AuthService) GenerateMagicLink(email string) (string, string, error) {
	secretKey := "your-very-secret-key" // This should be securely sourced, possibly from environment variables
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}

	link := "http://localhost:3000/auth/validate?token=" + tokenString
	return link, tokenString, nil
}

func (as *AuthService) ValidateToken(token string) (*Token, error) {
	// Retrieve the user and token details from the database
	user, err := as.getUserByToken(token)
	if err != nil {
		return nil, err
	}

	// Since sqlc handles the expiration check in SQL, if we get a user, the token is valid
	return &Token{
		ID:        user.MagicToken.String,
		UserID:    user.ID,
		ExpiresAt: user.TokenExpiration.Time,
	}, nil
}

func (ac *AuthService) getUserByToken(token string) (*db.User, error) {
	ctx := context.Background()
	user, err := ac.db.GetUserByToken(ctx, sql.NullString{String: token, Valid: true})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("token not found or expired")
		}
		return nil, err
	}
	return &user, nil
}
