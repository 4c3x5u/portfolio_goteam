package main

import (
	"database/sql"
	"net/http"
	"os"

	"server/api"
	boardAPI "server/api/board"
	columnAPI "server/api/column"
	loginAPI "server/api/login"
	registerAPI "server/api/register"
	taskAPI "server/api/task"
	"server/auth"
	boardTable "server/dbaccess/board"
	columnTable "server/dbaccess/column"
	taskTable "server/dbaccess/task"
	userTable "server/dbaccess/user"
	userboardTable "server/dbaccess/userboard"
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
	env := newEnv()
	if err := env.validate(); err != nil {
		log.Fatal(err.Error())
		os.Exit(2)
	}

	// Create dependencies that are shared by multiple handlers.
	db, err := sql.Open("postgres", env.DBConnStr)
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(3)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err.Error())
		os.Exit(4)
	}
	jwtGenerator := auth.NewJWTGenerator(env.JWTKey)
	userSelector := userTable.NewSelector(db)

	// Register handlers for API routes.
	mux := http.NewServeMux()

	mux.Handle("/register", registerAPI.NewHandler(
		registerAPI.NewValidator(
			registerAPI.NewUsernameValidator(),
			registerAPI.NewPasswordValidator(),
		),
		userSelector,
		registerAPI.NewPasswordHasher(),
		userTable.NewInserter(db),
		jwtGenerator,
		log,
	))

	mux.Handle("/login", loginAPI.NewHandler(
		loginAPI.NewValidator(),
		userSelector,
		loginAPI.NewPasswordComparator(),
		jwtGenerator,
		log,
	))

	mux.Handle("/board", api.NewHandler(
		auth.NewBearerTokenReader(),
		auth.NewJWTValidator(env.JWTKey),
		map[string]api.MethodHandler{
			http.MethodPost: boardAPI.NewPOSTHandler(
				boardAPI.NewNameValidator(),
				userboardTable.NewCounter(db),
				boardTable.NewInserter(db),
				log,
			),
			http.MethodDelete: boardAPI.NewDELETEHandler(
				boardAPI.NewIDValidator(),
				userboardTable.NewSelector(db),
				boardTable.NewDeleter(db),
				log,
			),
			http.MethodPatch: boardAPI.NewPATCHHandler(
				boardAPI.NewIDValidator(),
				boardAPI.NewNameValidator(),
				boardTable.NewSelector(db),
				userboardTable.NewSelector(db),
				boardTable.NewUpdater(db),
				log,
			),
		},
	))

	mux.Handle("/column", columnAPI.NewHandler(
		auth.NewBearerTokenReader(),
		auth.NewJWTValidator(env.JWTKey),
		columnAPI.NewIDValidator(),
		columnTable.NewSelector(db),
		userboardTable.NewSelector(db),
		columnTable.NewUpdater(db),
		log,
	))

	mux.Handle("/task", api.NewHandler(
		auth.NewBearerTokenReader(),
		auth.NewJWTValidator(jwtKey),
		map[string]api.MethodHandler{
			http.MethodPost: taskAPI.NewPOSTHandler(
				taskAPI.NewTitleValidator(),
				taskAPI.NewTitleValidator(),
				columnTable.NewSelector(db),
				userboardTable.NewSelector(db),
				taskTable.NewInserter(db),
				pkgLog.New(),
			),
		},
	))

	// Set up CORS.
	handler := midware.NewCORS(mux, env.ClientOrigin)

	// Serve the app using the ServeMux.
	log.Info("running server at port " + env.Port)
	if err := http.ListenAndServe(":"+env.Port, handler); err != nil {
		log.Fatal(err.Error())
		os.Exit(5)
	}
}
