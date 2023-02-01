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

// Server represents the goteam server as a whole and contains the top-level
// application logic.
type Server struct{ routeHandlers map[string]http.Handler }

// NewDefaultServer constructs and returns a Server with default dependencies.
func NewDefaultServer(
	dbPool *sql.DB,
	jwtKey string,
	logger log.Logger,
) Server {
	userSelector := db.NewUserSelector(dbPool)
	jwtGenerator := auth.NewJWTGenerator(jwtKey)
	return Server{
		routeHandlers: map[string]http.Handler{
			"/register": registerAPI.NewHandler(
				registerAPI.NewValidator(
					registerAPI.NewUsernameValidator(),
					registerAPI.NewPasswordValidator(),
				),
				userSelector,
				registerAPI.NewPasswordHasher(),
				db.NewUserInserter(dbPool),
				jwtGenerator,
				logger,
			),
			"/login": loginAPI.NewHandler(
				loginAPI.NewValidator(),
				userSelector,
				loginAPI.NewPasswordComparer(),
				jwtGenerator,
				logger,
			),
			"/board": boardAPI.NewHandler(
				auth.NewBearerTokenReader(),
				auth.NewJWTValidator(jwtKey),
				map[string]api.MethodHandler{
					http.MethodPost: boardAPI.NewPOSTHandler(
						boardAPI.NewPOSTValidator(),
						db.NewUserBoardCounter(dbPool),
						db.NewBoardInserter(dbPool),
						logger,
					),
					http.MethodDelete: boardAPI.NewDELETEHandler(
						boardAPI.NewDELETEValidator(),
						db.NewUserBoardSelector(dbPool),
						db.NewBoardDeleter(dbPool),
						logger,
					),
				},
			),
		},
	}
}

// Run serves the  API routes using the Server's route handlers.
func (s Server) Run(port string) error {
	mux := http.NewServeMux()
	for route, handler := range s.routeHandlers {
		mux.Handle(route, handler)
	}
	return http.ListenAndServe(":"+port, mux)
}

func main() {
	// Create a logger to be used throughout the app.
	logger := log.NewAppLogger()

	// Load environment variables from .env file.
	err := godotenv.Load()
	if err != nil {
		logger.Log(log.LevelFatal, err.Error())
		os.Exit(1)
	}

	// Create a database connection pool to be used throughout the app.
	dbPool, err := sql.Open("postgres", os.Getenv("DBCONNSTR"))
	if err != nil {
		logger.Log(log.LevelFatal, err.Error())
		os.Exit(1)
	}

	server := NewDefaultServer(dbPool, os.Getenv("JWTKEY"), logger)

	port := os.Getenv("PORT")
	logger.Log(log.LevelInfo, "running server at port "+port)
	if err := server.Run(port); err != nil {
		logger.Log(log.LevelFatal, err.Error())
		os.Exit(1)
	}
}
