package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/kxplxn/goteam/internal/api"
	boardAPI "github.com/kxplxn/goteam/internal/api/board"
	loginAPI "github.com/kxplxn/goteam/internal/api/login"
	registerAPI "github.com/kxplxn/goteam/internal/api/register"
	taskAPI "github.com/kxplxn/goteam/internal/api/task"
	tasksAPI "github.com/kxplxn/goteam/internal/api/tasks"
	teamAPI "github.com/kxplxn/goteam/internal/api/team"
	"github.com/kxplxn/goteam/pkg/auth"
	taskTAble "github.com/kxplxn/goteam/pkg/db/task"
	teamTable "github.com/kxplxn/goteam/pkg/db/team"
	userTable "github.com/kxplxn/goteam/pkg/db/user"
	legacyBoardTable "github.com/kxplxn/goteam/pkg/legacydb/board"
	legacyUserTable "github.com/kxplxn/goteam/pkg/legacydb/user"
	pkgLog "github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/pkg/token"
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
	env := api.NewEnv()
	if err := env.Validate(); err != nil {
		log.Fatal(err.Error())
		os.Exit(2)
	}

	cfgDynamoDB, err := config.LoadDefaultConfig(
		context.Background(), config.WithRegion(os.Getenv("AWS_REGION")),
	)
	if err != nil {
		log.Fatal(err.Error())
	}
	svcDynamo := dynamodb.NewFromConfig(cfgDynamoDB)

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
	jwtValidator := auth.NewJWTValidator(env.JWTKey)
	userSelector := legacyUserTable.NewSelector(db)

	// Register handlers for API routes.
	mux := http.NewServeMux()

	mux.Handle("/register", api.NewHandler(nil,
		map[string]api.MethodHandler{
			http.MethodPost: registerAPI.NewPostHandler(
				registerAPI.NewUserValidator(
					registerAPI.NewUsernameValidator(),
					registerAPI.NewPasswordValidator(),
				),
				token.DecodeInvite,
				registerAPI.NewPasswordHasher(),
				userTable.NewInserter(svcDynamo),
				token.EncodeAuth,
				log,
			),
		},
	))

	mux.Handle("/login", api.NewHandler(nil,
		map[string]api.MethodHandler{
			http.MethodPost: loginAPI.NewPostHandler(
				loginAPI.NewValidator(),
				userTable.NewRetriever(svcDynamo),
				loginAPI.NewPasswordComparator(),
				token.EncodeAuth,
				log,
			),
		},
	))

	mux.Handle("/team", api.NewHandler(nil,
		map[string]api.MethodHandler{
			http.MethodGet: teamAPI.NewGetHandler(
				token.DecodeAuth,
				teamTable.NewRetriever(svcDynamo),
				teamTable.NewInserter(svcDynamo),
				log,
			),
		},
	))

	boardIDValidator := boardAPI.NewIDValidator()
	boardNameValidator := boardAPI.NewNameValidator()
	boardSelector := legacyBoardTable.NewSelector(db)
	boardInserter := legacyBoardTable.NewInserter(db)
	mux.Handle("/board", api.NewHandler(
		jwtValidator,
		map[string]api.MethodHandler{
			http.MethodPost: boardAPI.NewPOSTHandler(
				userSelector,
				boardNameValidator,
				legacyBoardTable.NewCounter(db),
				boardInserter,
				log,
			),
			http.MethodDelete: boardAPI.NewDeleteHandler(
				token.DecodeAuth,
				token.DecodeState,
				teamTable.NewBoardDeleter(svcDynamo, svcDynamo),
				log,
			),
			http.MethodPatch: boardAPI.NewPATCHHandler(
				userSelector,
				boardIDValidator,
				boardNameValidator,
				boardSelector,
				legacyBoardTable.NewUpdater(db),
				log,
			),
		},
	))

	mux.Handle("/tasks", api.NewHandler(
		jwtValidator,
		map[string]api.MethodHandler{
			http.MethodPatch: tasksAPI.NewPatchHandler(
				token.DecodeAuth,
				token.DecodeState,
				tasksAPI.NewColNoValidator(),
				taskTAble.NewMultiUpdater(svcDynamo),
				token.EncodeState,
				log,
			),
			http.MethodGet: tasksAPI.NewGetHandler(
				token.DecodeAuth,
				taskTAble.NewMultiRetriever(svcDynamo),
				log,
			),
		},
	))

	taskTitleValidator := taskAPI.NewTitleValidator()
	mux.Handle("/task", api.NewHandler(
		jwtValidator,
		map[string]api.MethodHandler{
			http.MethodPost: taskAPI.NewPostHandler(
				token.DecodeAuth,
				token.DecodeState,
				taskTitleValidator,
				taskTitleValidator,
				taskAPI.NewColNoValidator(),
				taskTAble.NewInserter(svcDynamo),
				token.EncodeState,
				log,
			),
			http.MethodPatch: taskAPI.NewPatchHandler(
				token.DecodeAuth,
				token.DecodeState,
				taskTitleValidator,
				taskTitleValidator,
				taskTAble.NewUpdater(svcDynamo),
				log,
			),
			http.MethodDelete: taskAPI.NewDeleteHandler(
				token.DecodeAuth,
				token.DecodeState,
				taskTAble.NewDeleter(svcDynamo),
				token.EncodeState,
				log,
			),
		},
	))

	// Serve the app using the ServeMux.
	log.Info("running server at port " + env.Port)
	if err := http.ListenAndServe(":"+env.Port, mux); err != nil {
		log.Fatal(err.Error())
		os.Exit(5)
	}
}
