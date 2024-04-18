package controllers

import (
	"net/http"
	"resons/api/api/db"
	"resons/api/api/services"
	"resons/api/api/types"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"

	"time"
)

type AuthController struct {
	EmailService *services.EmailService
	AuthService  *services.AuthService
}

func NewAuthController(queries *db.Queries) *AuthController {
	emailService := services.NewEmailService("localhost", 1025, "", "", "contact@vibespot.com")
	authService := services.NewAuthService(queries)

	return &AuthController{
		EmailService: emailService,
		AuthService:  authService,
	}
}

func (ac *AuthController) HandleVerify(c echo.Context) error {
	token := c.QueryParam("token")

	// Validate the token and check expiration
	userToken, _ := ac.AuthService.ValidateToken(token)
	userID := userToken.UserID

	// Create JWT token
	claims := types.JWTCustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := jwtToken.SignedString([]byte("secret"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to sign token")
	}

	c.Response().Header().Set(echo.HeaderAuthorization, "Bearer "+tokenString)

	return c.String(http.StatusOK, "User "+userToken.UserID.String()+" verified successfully")
}

func (ac *AuthController) HandleLogin(c echo.Context) error {
	email := c.FormValue("email")
	link, tokenString, err := ac.AuthService.GenerateMagicLink(email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to generate magic link"})
	}

	ac.AuthService.UpsertUserWithToken(c, tokenString, email)

	err = ac.EmailService.SendMagicLink(email, link)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to send magic link"})
	}

	return c.String(http.StatusOK, "Magic link sent successfully")
}
