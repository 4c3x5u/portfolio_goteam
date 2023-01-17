package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	boardAPI "server/api/board"
	loginAPI "server/api/login"
	registerAPI "server/api/register"
	"server/auth"
	"server/db"
)

func main() {
	// Create dependencies that are shared by multiple handlers.
	conn, err := sql.Open("postgres", os.Getenv("DBCONNSTR"))
	if err != nil {
		log.Fatal(err)
	}
	connCloser := db.NewConnCloser(conn)
	jwtKey := os.Getenv("JWTKEY")
	jwtGenerator := auth.NewJWTGenerator(jwtKey)
	userSelector := db.NewUserSelector(conn)

	// Register handlers for API endpoints.
	mux := http.NewServeMux()
	mux.Handle("/register", registerAPI.NewHandler(
		registerAPI.NewRequestValidator(
			registerAPI.NewUsernameValidator(),
			registerAPI.NewPasswordValidator(),
		),
		userSelector,
		registerAPI.NewPasswordHasher(),
		db.NewUserInserter(conn),
		jwtGenerator,
		connCloser,
	))
	mux.Handle("/login", loginAPI.NewHandler(
		userSelector,
		loginAPI.NewPasswordComparer(),
		jwtGenerator,
		connCloser,
	))
	mux.Handle("/board", boardAPI.NewHandler(
		auth.NewJWTValidator(jwtKey),
	))

	// Serve the app using the ServeMux.
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
