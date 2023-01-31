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
	logger := log.NewAppLogger()

	// Load environment variables from .env file.
	err := godotenv.Load()
	if err != nil {
		logger.Log(log.LevelFatal, err.Error())
		os.Exit(1)
	}

	// Create dependencies that are shared by multiple handlers.
	conn, err := sql.Open("postgres", os.Getenv("DBCONNSTR"))
	connCloser := db.NewConnCloser(conn)
	defer connCloser.Close()
	if err != nil {
		logger.Log(log.LevelFatal, err.Error())
		os.Exit(1)
	}

	jwtKey := os.Getenv("JWTKEY")
	jwtGenerator := auth.NewJWTGenerator(jwtKey)

	userSelector := db.NewUserSelector(conn)

	// Register handlers for API routes.
	mux := http.NewServeMux()

	mux.Handle("/register", registerAPI.NewHandler(
		registerAPI.NewValidator(
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
		loginAPI.NewValidator(),
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
			http.MethodPost: boardAPI.NewPOSTHandler(
				boardAPI.NewPOSTValidator(),
				db.NewUserBoardCounter(conn),
				db.NewBoardInserter(conn),
				logger,
			),
			http.MethodDelete: boardAPI.NewDELETEHandler(
				boardAPI.NewDELETEValidator(),
				db.NewUserBoardSelector(conn),
				db.NewBoardDeleter(conn),
				logger,
			),
		},
	))

	// Serve the app using the ServeMux.
	port := os.Getenv("PORT")
	logger.Log(log.LevelInfo, "running server at port "+port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		logger.Log(log.LevelFatal, err.Error())
		os.Exit(1)
	}
}
