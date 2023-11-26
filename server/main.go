package main

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/kxplxn/goteam/server/api"
	boardAPI "github.com/kxplxn/goteam/server/api/board"
	columnAPI "github.com/kxplxn/goteam/server/api/column"
	loginAPI "github.com/kxplxn/goteam/server/api/login"
	registerAPI "github.com/kxplxn/goteam/server/api/register"
	subtaskAPI "github.com/kxplxn/goteam/server/api/subtask"
	taskAPI "github.com/kxplxn/goteam/server/api/task"
	"github.com/kxplxn/goteam/server/auth"
	boardTable "github.com/kxplxn/goteam/server/dbaccess/board"
	columnTable "github.com/kxplxn/goteam/server/dbaccess/column"
	subtaskTable "github.com/kxplxn/goteam/server/dbaccess/subtask"
	taskTable "github.com/kxplxn/goteam/server/dbaccess/task"
	teamTable "github.com/kxplxn/goteam/server/dbaccess/team"
	userTable "github.com/kxplxn/goteam/server/dbaccess/user"
	pkgLog "github.com/kxplxn/goteam/server/log"
	"github.com/kxplxn/goteam/server/midware"
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

	// Register handlers for API routes.
	mux := http.NewServeMux()

	teamSelectorByInvCode := teamTable.NewSelectorByInvCode(db)
	mux.Handle("/register", registerAPI.NewHandler(
		registerAPI.NewUserValidator(
			registerAPI.NewUsernameValidator(),
			registerAPI.NewPasswordValidator(),
		),
		registerAPI.NewInviteCodeValidator(),
		teamSelectorByInvCode,
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
	boardSelector := boardTable.NewSelector(db)
	mux.Handle("/board", api.NewHandler(
		bearerTokenReader,
		jwtValidator,
		map[string]api.MethodHandler{
			http.MethodGet: boardAPI.NewGETHandler(
				userSelector,
				boardIDValidator,
				boardTable.NewRecursiveSelector(db),
				teamTable.NewSelector(db),
				userTable.NewSelectorByTeamID(db),
				boardTable.NewSelectorByTeamID(db),
				log,
			),
			http.MethodPost: boardAPI.NewPOSTHandler(
				userSelector,
				boardNameValidator,
				boardTable.NewCounter(db),
				boardTable.NewInserter(db),
				log,
			),
			http.MethodDelete: boardAPI.NewDELETEHandler(
				userSelector,
				boardIDValidator,
				boardSelector,
				boardTable.NewDeleter(db),
				log,
			),
			http.MethodPatch: boardAPI.NewPATCHHandler(
				userSelector,
				boardIDValidator,
				boardNameValidator,
				boardSelector,
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
				userSelector,
				columnAPI.NewIDValidator(),
				columnSelector,
				boardSelector,
				columnTable.NewUpdater(db),
				log,
			),
		},
	))

	taskIDValidator := taskAPI.NewIDValidator()
	taskTitleValidator := taskAPI.NewTitleValidator()
	taskSelector := taskTable.NewSelector(db)
	mux.Handle("/task", api.NewHandler(
		bearerTokenReader,
		jwtValidator,
		map[string]api.MethodHandler{
			http.MethodPost: taskAPI.NewPOSTHandler(
				userSelector,
				taskTitleValidator,
				taskTitleValidator,
				columnSelector,
				boardSelector,
				taskTable.NewInserter(db),
				log,
			),
			http.MethodPatch: taskAPI.NewPATCHHandler(
				userSelector,
				taskIDValidator,
				taskTitleValidator,
				taskTitleValidator,
				taskSelector,
				columnSelector,
				boardSelector,
				taskTable.NewUpdater(db),
				log,
			),
			http.MethodDelete: taskAPI.NewDELETEHandler(
				userSelector,
				taskIDValidator,
				taskSelector,
				columnSelector,
				boardSelector,
				taskTable.NewDeleter(db),
				log,
			),
		},
	))

	mux.Handle("/subtask", api.NewHandler(
		bearerTokenReader,
		jwtValidator,
		map[string]api.MethodHandler{
			http.MethodPatch: subtaskAPI.NewPATCHHandler(
				userSelector,
				subtaskAPI.NewIDValidator(),
				subtaskTable.NewSelector(db),
				taskSelector,
				columnSelector,
				boardSelector,
				subtaskTable.NewUpdater(db),
				pkgLog.New(),
			),
		},
	))

	// Set up access control headers required by the endpoints.
	handler := midware.NewAccessControl(mux, env.ClientOrigin)

	// Serve the app using the ServeMux.
	log.Info("running server at port " + env.Port)
	if err := http.ListenAndServe(":"+env.Port, handler); err != nil {
		log.Fatal(err.Error())
		os.Exit(5)
	}
}
