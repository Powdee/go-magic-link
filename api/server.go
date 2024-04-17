package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"resons/v0/api/api/controllers"
	"resons/v0/api/api/db"
	"resons/v0/api/api/middleware"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func initDB() *db.Queries {
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

	return db.New(database)
}

func main() {
	e := echo.New()
	queries := initDB()

	authController := controllers.NewAuthController(queries)
	jwtMiddleware := middleware.JWTMiddleware("your-very-secret-key")

	e.POST("/auth/login", authController.HandleLogin)
	e.GET("/auth/verify", authController.HandleVerify)

	e.POST("/events/user/:userId/upload", func(c echo.Context) error { return c.String(http.StatusOK, "authenticated") }, jwtMiddleware)

	e.Logger.Fatal(e.Start(":8080"))
}
