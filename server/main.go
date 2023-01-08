package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"server/token"

	loginAPI "server/api/login"
	registerAPI "server/api/register"
	"server/db"

	"github.com/golang-jwt/jwt/v4"
)

func main() {
	// Create dependencies that are shared by multiple handlers.
	conn, err := sql.Open("postgres", os.Getenv("DBCONNSTR"))
	if err != nil {
		log.Fatal(err)
	}

	jwtGenerator := token.NewJWTGenerator(
		os.Getenv("JWTSIGNATURE"),
		jwt.SigningMethodHS256,
	)

	// Register handlers for API endpoints.
	mux := http.NewServeMux()

	mux.Handle("/register", registerAPI.NewHandler(
		registerAPI.NewRequestValidator(
			registerAPI.NewUsernameValidator(),
			registerAPI.NewPasswordValidator(),
		),
		db.NewUserReader(conn),
		registerAPI.NewPasswordHasher(),
		db.NewUserCreator(conn),
		jwtGenerator,
	))

	mux.Handle("/login", loginAPI.NewHandler(
		db.NewUserReader(conn),
		loginAPI.NewPasswordComparer(),
		jwtGenerator,
	))

	// Serve the API routes using the ServeMux.
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
