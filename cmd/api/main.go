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
	lgcBoardAPI "github.com/kxplxn/goteam/internal/api/board"
	tasksAPI "github.com/kxplxn/goteam/internal/api/tasks"
	taskAPI "github.com/kxplxn/goteam/internal/api/tasks/task"
	teamAPI "github.com/kxplxn/goteam/internal/api/team"
	boardAPI "github.com/kxplxn/goteam/internal/api/team/board"
	loginAPI "github.com/kxplxn/goteam/internal/api/user/login"
	registerAPI "github.com/kxplxn/goteam/internal/api/user/register"
	"github.com/kxplxn/goteam/pkg/auth"
	taskTable "github.com/kxplxn/goteam/pkg/db/task"
	teamTable "github.com/kxplxn/goteam/pkg/db/team"
	userTable "github.com/kxplxn/goteam/pkg/db/user"
	lgcBoardTable "github.com/kxplxn/goteam/pkg/legacydb/board"
	lgcUserTable "github.com/kxplxn/goteam/pkg/legacydb/user"
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
	userSelector := lgcUserTable.NewSelector(db)

	// Register handlers for API routes.
	mux := http.NewServeMux()

	mux.Handle("/user/register", api.NewHandler(
		nil,
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

	mux.Handle("/user/login", api.NewHandler(nil, map[string]api.MethodHandler{
		http.MethodPost: loginAPI.NewPostHandler(
			loginAPI.NewValidator(),
			userTable.NewRetriever(svcDynamo),
			loginAPI.NewPasswordComparator(),
			token.EncodeAuth,
			log,
		),
	}))

	mux.Handle("/team", api.NewHandler(nil, map[string]api.MethodHandler{
		http.MethodGet: teamAPI.NewGetHandler(
			token.DecodeAuth,
			teamTable.NewRetriever(svcDynamo),
			teamTable.NewInserter(svcDynamo),
			log,
		),
	}))

	mux.Handle("/team/board", api.NewHandler(nil, map[string]api.MethodHandler{
		http.MethodDelete: boardAPI.NewDeleteHandler(
			token.DecodeAuth,
			token.DecodeState,
			teamTable.NewBoardDeleter(svcDynamo),
			log,
		),
	}))

	// TODO: remove once fully migrated to DynamoDB
	boardNameValidator := lgcBoardAPI.NewNameValidator()
	mux.Handle("/board", api.NewHandler(
		jwtValidator,
		map[string]api.MethodHandler{
			http.MethodPost: lgcBoardAPI.NewPOSTHandler(
				userSelector,
				boardNameValidator,
				lgcBoardTable.NewCounter(db),
				lgcBoardTable.NewInserter(db),
				log,
			),
			http.MethodPatch: lgcBoardAPI.NewPatchHandler(
				token.DecodeAuth,
				token.DecodeState,
				lgcBoardAPI.NewIDValidator(),
				boardNameValidator,
				teamTable.NewBoardUpdater(svcDynamo),
				log,
			),
		},
	))

	taskTitleValidator := taskAPI.NewTitleValidator()
	mux.Handle("/tasks", api.NewHandler(
		jwtValidator,
		map[string]api.MethodHandler{
			http.MethodPatch: tasksAPI.NewPatchHandler(
				token.DecodeAuth,
				token.DecodeState,
				tasksAPI.NewColNoValidator(),
				taskTable.NewMultiUpdater(svcDynamo),
				token.EncodeState,
				log,
			),
			http.MethodGet: tasksAPI.NewGetHandler(
				token.DecodeAuth,
				taskTable.NewMultiRetriever(svcDynamo),
				log,
			),
		},
	))

	mux.Handle("/tasks/task", api.NewHandler(
		jwtValidator,
		map[string]api.MethodHandler{
			http.MethodPost: taskAPI.NewPostHandler(
				token.DecodeAuth,
				token.DecodeState,
				taskTitleValidator,
				taskTitleValidator,
				taskAPI.NewColNoValidator(),
				taskTable.NewInserter(svcDynamo),
				token.EncodeState,
				log,
			),
			http.MethodPatch: taskAPI.NewPatchHandler(
				token.DecodeAuth,
				token.DecodeState,
				taskTitleValidator,
				taskTitleValidator,
				taskTable.NewUpdater(svcDynamo),
				log,
			),
			http.MethodDelete: taskAPI.NewDeleteHandler(
				token.DecodeAuth,
				token.DecodeState,
				taskTable.NewDeleter(svcDynamo),
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
