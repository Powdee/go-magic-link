package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"resons/v0/api/api/db"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"

	"gopkg.in/gomail.v2"

	"errors"
	"time"
)

type Token struct {
	ID        string
	UserID    string    // The ID of the user who owns this token
	ExpiresAt time.Time // Expiration time of the token
}

func sendMagicLink(email string, link string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "contact@vibespot.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Your Magic Login Link")
	m.SetBody("text/html", "Click <a href=\""+link+"\">here</a> to log in.")

	d := gomail.NewDialer("localhost", 1025, "", "")

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func getUserByToken(token string) (*db.User, error) {
	ctx := context.Background()
	user, err := queries.GetUserByToken(ctx, sql.NullString{String: token, Valid: true})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("token not found or expired")
		}
		return nil, err
	}
	return &user, nil
}

func validateToken(token string) (*Token, error) {
	// Retrieve the user and token details from the database
	user, err := getUserByToken(token)
	if err != nil {
		return nil, err
	}

	// Since sqlc handles the expiration check in SQL, if we get a user, the token is valid
	return &Token{
		ID:        user.MagicToken.String,     // Access the value of user.MagicToken using .String method
		UserID:    strconv.Itoa(int(user.ID)), // Assuming ID is an int and needs to be a string
		ExpiresAt: user.TokenExpiration.Time,  // Assuming TokenExpiration is sql.NullTime
	}, nil
}

var queries *db.Queries

func initDB() {
	var err error
	// todo: use env variables
	username := "erikkurjak" // os.Getenv("DB_USER")
	password := "w!ndow11"   // os.Getenv("DB_PASSWORD")
	hostname := "localhost"  // os.Getenv("DB_HOST")
	dbname := "momenify"     // os.Getenv("DB_NAME")
	fmt.Println(username, password, hostname, dbname)
	// Create the connection string
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", username, password, hostname, dbname)

	// Open the database connection
	database, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}

	// Check the database connection
	if err = database.Ping(); err != nil {
		log.Fatal("Failed to connect to DB: ", err)
	}

	// Initialize the generated Queries struct
	queries = db.New(database)
}

func generateMagicLink(email string) (string, string, error) {
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

type jwtCustomClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func main() {
	e := echo.New()
	initDB()

	jwtMiddleware := echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte("your-very-secret-key"),
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			log.Printf("c: %v", c)

			return new(jwtCustomClaims)
		},
		TokenLookup: "header:Authorization",
	})

	e.POST("/auth/login", func(c echo.Context) error {
		email := c.FormValue("email")
		link, tokenString, err := generateMagicLink(email)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to generate magic link"})
		}
		log.Printf("c: %v", c)

		expiration := time.Now().Add(24 * time.Hour)
		magicToken := sql.NullString{String: tokenString, Valid: true}
		expirationTime := sql.NullTime{Time: expiration, Valid: true}

		if err := queries.UpsertUserWithToken(context.Background(), db.UpsertUserWithTokenParams{
			Email:           email,
			MagicToken:      magicToken,
			TokenExpiration: expirationTime,
		}); err != nil {
			log.Printf("Database error: %v", err)
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Unable to upsert user"})
		}

		err = sendMagicLink(email, link)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to send magic link"})
		}

		return c.String(http.StatusOK, "Magic link sent successfully")
	})

	e.GET("/auth/callback", func(c echo.Context) error {
		token := c.QueryParam("token")

		// Validate the token and check expiration
		userToken, _ := validateToken(token)
		userID, err := strconv.Atoi(userToken.UserID) // Convert UserID string back to int, handle error if necessary
		if err != nil {
			return c.Redirect(http.StatusTemporaryRedirect, "/login?error="+err.Error())
		}

		// Create JWT token
		claims := jwtCustomClaims{
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
		// Redirect the user to their dashboard
		return c.String(http.StatusOK, "User "+userToken.UserID+" verified successfully")
	})

	e.POST("/events/user/:userId/upload", func(c echo.Context) error { return c.String(http.StatusOK, "authenticated") }, jwtMiddleware)

	e.Logger.Fatal(e.Start(":8080"))
}
