package main

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/kxplxn/goteam/server/api"
	boardAPI "github.com/kxplxn/goteam/server/api/board"
	columnAPI "github.com/kxplxn/goteam/server/api/column"
	loginAPI "github.com/kxplxn/goteam/server/api/login"
	registerAPI "github.com/kxplxn/goteam/server/api/register"
	taskAPI "github.com/kxplxn/goteam/server/api/task"
	"github.com/kxplxn/goteam/server/auth"
	boardTable "github.com/kxplxn/goteam/server/dbaccess/board"
	columnTable "github.com/kxplxn/goteam/server/dbaccess/column"
	taskTable "github.com/kxplxn/goteam/server/dbaccess/task"
	userTable "github.com/kxplxn/goteam/server/dbaccess/user"
	userboardTable "github.com/kxplxn/goteam/server/dbaccess/userboard"
	pkgLog "github.com/kxplxn/goteam/server/log"
	"github.com/kxplxn/goteam/server/midware"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	// Create a logger for the app.
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

	// Create dependencies that are used by multiple handlers.
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
	bearerTokenReader := auth.NewBearerTokenReader()
	jwtValidator := auth.NewJWTValidator(env.JWTKey)
	userSelector := userTable.NewSelector(db)
	columnSelector := columnTable.NewSelector(db)
	userBoardSelector := userboardTable.NewSelector(db)

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

	boardIDValidator := boardAPI.NewIDValidator()
	boardNameValidator := boardAPI.NewNameValidator()
	mux.Handle("/board", api.NewHandler(
		bearerTokenReader,
		jwtValidator,
		map[string]api.MethodHandler{
			http.MethodPost: boardAPI.NewPOSTHandler(
				boardNameValidator,
				userboardTable.NewCounter(db),
				boardTable.NewInserter(db),
				log,
			),
			http.MethodDelete: boardAPI.NewDELETEHandler(
				boardIDValidator,
				userBoardSelector,
				boardTable.NewDeleter(db),
				log,
			),
			http.MethodPatch: boardAPI.NewPATCHHandler(
				boardIDValidator,
				boardNameValidator,
				boardTable.NewSelector(db),
				userBoardSelector,
				boardTable.NewUpdater(db),
				log,
			),
		},
	))

	mux.Handle("/column", api.NewHandler(
		bearerTokenReader,
		jwtValidator,
		map[string]api.MethodHandler{
			http.MethodPatch: columnAPI.NewPATCHHandler(
				columnAPI.NewIDValidator(),
				columnSelector,
				userBoardSelector,
				columnTable.NewUpdater(db),
				log,
			),
		},
	))

	mux.Handle("/task", api.NewHandler(
		bearerTokenReader,
		jwtValidator,
		map[string]api.MethodHandler{
			http.MethodPost: taskAPI.NewPOSTHandler(
				taskAPI.NewTitleValidator(),
				taskAPI.NewTitleValidator(),
				columnSelector,
				userBoardSelector,
				taskTable.NewInserter(db),
				log,
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
