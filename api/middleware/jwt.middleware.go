package middleware

import (
	"resons/v0/api/api/types"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware(secretKey string) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(secretKey),
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(types.JWTCustomClaims)
		},
		TokenLookup: "header:Authorization",
	})
}
