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
	"server/midware"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

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
	env := map[string]string{
		Port:         os.Getenv(Port),
		DBConnStr:    os.Getenv(DBConnStr),
		JWTKey:       os.Getenv(JWTKey),
		ClientOrigin: os.Getenv(ClientOrigin),
	}
	for name, value := range env {
		if value == "" {
			log.Fatal(name + " env var was empty")
			os.Exit(2)
		}
	}

	// Create dependencies that are shared by multiple handlers.
	dbConn, err := sql.Open("postgres", env[DBConnStr])
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(3)
	}
	if err = dbConn.Ping(); err != nil {
		log.Fatal(err.Error())
		os.Exit(4)
	}
	jwtGenerator := auth.NewJWTGenerator(env[JWTKey])
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
		auth.NewJWTValidator(env[JWTKey]),
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

	// Set up CORS.
	handler := midware.NewCORS(mux, env[ClientOrigin])

	// Serve the app using the ServeMux.
	log.Info("running server at port " + env[Port])
	if err := http.ListenAndServe(":"+env[Port], handler); err != nil {
		log.Fatal(err.Error())
		os.Exit(5)
	}
}

const (
	// Port is the name of the environment variable used for deciding what port
	// to run the server on.
	Port = "PORT"

	// DBConnStr is the name of the environment variable used for connecting to
	// the database.
	DBConnStr = "DBCONNSTR"

	// JWTKey is the name of the environment variable used for signing JWTs.
	JWTKey = "JWTKEY"

	// ClientOrigin is the name of the environment variable used to set up CORS
	// with the client app.
	ClientOrigin = "CLIENTORIGIN"
)
