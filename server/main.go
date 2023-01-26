package main

import (
	"database/sql"
	"net/http"
	"os"

	"server/api"
	boardAPI "server/api/board"
	loginAPI "server/api/login"
	registerAPI "server/api/register"
	"server/auth"
	"server/db"
	"server/log"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	logger := log.NewBasicLogger()

	// Load environment variables from .env file.
	err := godotenv.Load()
	if err != nil {
		logger.Log(log.LevelFatal, err.Error())
	}

	// Create dependencies that are shared by multiple handlers.
	conn, err := sql.Open("postgres", os.Getenv("DBCONNSTR"))
	connCloser := db.NewConnCloser(conn)
	defer connCloser.Close()
	if err != nil {
		logger.Log(log.LevelFatal, err.Error())
	}

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
		logger,
	))

	mux.Handle("/login", loginAPI.NewHandler(
		userSelector,
		loginAPI.NewPasswordComparer(),
		jwtGenerator,
		connCloser,
		logger,
	))

	mux.Handle("/board", boardAPI.NewHandler(
		auth.NewBearerTokenReader(),
		auth.NewJWTValidator(jwtKey),
		map[string]api.MethodHandler{
			http.MethodPost: boardAPI.NewPostHandler(
				db.NewUserBoardCounter(conn),
				db.NewBoardInserter(conn),
				logger,
			),
			http.MethodDelete: boardAPI.NewDeleteHandler(
				db.NewUserBoardSelector(conn),
				db.NewBoardDeleter(conn),
				logger,
			),
		},
	))

	// Serve the app using the ServeMux.
	if err := http.ListenAndServe(":8080", mux); err != nil {
		logger.Log(log.LevelFatal, err.Error())
		os.Exit(1)
	}
}
