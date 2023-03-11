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
	pkgLog "server/log"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

// envPORT is the name of the environment variable used for deciding what port
// to run the server on.
const envPORT = "PORT"

// envDBCONNSTR is the name of the environment variable used for connecting to
// the database.
const envDBCONNSTR = "DBCONNSTR"

// envJWTKEY is the name of the environment variable used for signing JWTs.
const envJWTKEY = "JWTKEY"

func main() {
	// Create a log for the app.
	log := pkgLog.New()

	// Load environment variables from .env file.
	err := godotenv.Load()
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	// Ensure that the necessary env vars were set.
	port := os.Getenv(envPORT)
	dbConnStr := os.Getenv(envDBCONNSTR)
	jwtKey := os.Getenv(envJWTKEY)
	for name, value := range map[string]string{
		envPORT:      port,
		envDBCONNSTR: dbConnStr,
		envJWTKEY:    jwtKey,
	} {
		if value == "" {
			log.Fatal(name + " env var was empty")
			os.Exit(1)
		}
	}

	// Create dependencies that are shared by multiple handlers.
	dbConn, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}
	if err = dbConn.Ping(); err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}
	jwtGenerator := auth.NewJWTGenerator(jwtKey)
	userSelector := db.NewUserSelector(dbConn)

	// Register handlers for API routes.
	mux := http.NewServeMux()

	mux.Handle("/register", registerAPI.NewHandler(
		registerAPI.NewValidator(
			registerAPI.NewUsernameValidator(),
			registerAPI.NewPasswordValidator(),
		),
		userSelector,
		registerAPI.NewPasswordHasher(),
		db.NewUserInserter(dbConn),
		jwtGenerator,
		log,
	))

	mux.Handle("/login", loginAPI.NewHandler(
		loginAPI.NewValidator(),
		userSelector,
		loginAPI.NewPasswordComparer(),
		jwtGenerator,
		log,
	))

	mux.Handle("/board", boardAPI.NewHandler(
		auth.NewBearerTokenReader(),
		auth.NewJWTValidator(jwtKey),
		map[string]api.MethodHandler{
			http.MethodPost: boardAPI.NewPOSTHandler(
				boardAPI.NewPOSTValidator(),
				db.NewUserBoardCounter(dbConn),
				db.NewBoardInserter(dbConn),
				log,
			),
			http.MethodDelete: boardAPI.NewDELETEHandler(
				boardAPI.NewDELETEValidator(),
				db.NewUserBoardSelector(dbConn),
				db.NewBoardDeleter(dbConn),
				log,
			),
		},
	))

	// Serve the app using the ServeMux.
	log.Info("running server at port " + port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}
}
