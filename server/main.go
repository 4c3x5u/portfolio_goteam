package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	loginAPI "server/api/login"
	registerAPI "server/api/register"
	"server/cookie"
	"server/db"
)

func main() {
	// Create dependencies that are shared by multiple handlers.
	conn, err := sql.Open("postgres", os.Getenv("DBCONNSTR"))
	if err != nil {
		log.Fatal(err)
	}
	connCloser := db.NewConnCloser(conn)
	jwtGenerator := cookie.NewJWTGenerator(os.Getenv("JWTSIGNATURE"))

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
		connCloser,
	))
	mux.Handle("/login", loginAPI.NewHandler(
		db.NewUserReader(conn),
		loginAPI.NewPasswordComparer(),
		jwtGenerator,
		connCloser,
	))

	// Serve the app using the ServeMux.
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
