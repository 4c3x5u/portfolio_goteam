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
type Server struct {
	port          string
	logger        log.Logger
	routeHandlers map[string]http.Handler
}

// NewServer constructs and returns a Server with default dependencies.
func NewServer(logger log.Logger) Server {
	// Get environment variables and ensure they're not empty.
	dbConnStr := os.Getenv("DBCONNSTR")
	jwtKey := os.Getenv("JWTKEY")
	port := os.Getenv("PORT")
	for description, value := range map[string]string{
		"database connection string": dbConnStr,
		"jwt signature key":          jwtKey,
		"port variable":              port,
	} {
		if value == "" {
			logger.Log(log.LevelFatal, description+" was empty")
			os.Exit(1)
		}
	}

	// Create dependencies needed by multiple route handlers.
	dbPool, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		logger.Log(log.LevelFatal, err.Error())
		os.Exit(1)
	}
	userSelector := db.NewUserSelector(dbPool)
	jwtGenerator := auth.NewJWTGenerator(jwtKey)

	// Return the server with the port to run and route handlers.
	return Server{
		port: port,
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

// Run serves API routes using the Server's route handlers.
func (s Server) Run() error {
	s.logger.Log(log.LevelInfo, "running server at port "+s.port)
	mux := http.NewServeMux()
	for route, handler := range s.routeHandlers {
		mux.Handle(route, handler)
	}
	return http.ListenAndServe(":"+s.port, mux)
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

	// Create a new server and run it.
	if err := NewServer(logger).Run(); err != nil {
		logger.Log(log.LevelFatal, err.Error())
		os.Exit(1)
	}
}
