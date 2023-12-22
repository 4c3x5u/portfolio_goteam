package main

import (
	"context"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/joho/godotenv"

	taskAPI "github.com/kxplxn/goteam/internal/task/task"
	tasksAPI "github.com/kxplxn/goteam/internal/task/tasks"
	teamAPI "github.com/kxplxn/goteam/internal/team"
	boardAPI "github.com/kxplxn/goteam/internal/team/board"
	loginAPI "github.com/kxplxn/goteam/internal/user/login"
	registerAPI "github.com/kxplxn/goteam/internal/user/register"
	"github.com/kxplxn/goteam/pkg/api"
	"github.com/kxplxn/goteam/pkg/db/tasktable"
	"github.com/kxplxn/goteam/pkg/db/teamtable"
	"github.com/kxplxn/goteam/pkg/db/usertable"
	pkgLog "github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/pkg/token"
)

func main() {
	// create a logger
	log := pkgLog.New()

	// load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	// ensure that the necessary environment vars were set
	env := api.NewEnv()
	if err := env.Validate(); err != nil {
		log.Fatal(err.Error())
		os.Exit(2)
	}

	// create DynamoDB client
	dbCfg, err := config.LoadDefaultConfig(
		context.Background(), config.WithRegion(os.Getenv("AWS_REGION")),
	)
	if err != nil {
		log.Fatal(err.Error())
	}
	db := dynamodb.NewFromConfig(dbCfg)

	// register handlers for HTTP routes
	mux := http.NewServeMux()

	mux.Handle("/user/register", api.NewHandler(map[string]api.MethodHandler{
		http.MethodPost: registerAPI.NewPostHandler(
			registerAPI.NewUserValidator(
				registerAPI.NewUsernameValidator(),
				registerAPI.NewPasswordValidator(),
			),
			token.DecodeInvite,
			registerAPI.NewPasswordHasher(),
			usertable.NewInserter(db),
			token.EncodeAuth,
			log,
		),
	}))

	mux.Handle("/user/login", api.NewHandler(map[string]api.MethodHandler{
		http.MethodPost: loginAPI.NewPostHandler(
			loginAPI.NewValidator(),
			usertable.NewRetriever(db),
			loginAPI.NewPasswordComparator(),
			token.EncodeAuth,
			log,
		),
	}))

	mux.Handle("/team", api.NewHandler(map[string]api.MethodHandler{
		http.MethodGet: teamAPI.NewGetHandler(
			token.DecodeAuth,
			teamtable.NewRetriever(db),
			teamtable.NewInserter(db),
			log,
		),
	}))

	mux.Handle("/team/board", api.NewHandler(map[string]api.MethodHandler{
		http.MethodPost: boardAPI.NewPostHandler(
			token.DecodeAuth,
			token.DecodeState,
			boardAPI.NewNameValidator(),
			teamtable.NewBoardInserter(db),
			token.EncodeState,
			log,
		),
		http.MethodPatch: boardAPI.NewPatchHandler(
			token.DecodeAuth,
			token.DecodeState,
			boardAPI.NewIDValidator(),
			boardAPI.NewNameValidator(),
			teamtable.NewBoardUpdater(db),
			log,
		),
		http.MethodDelete: boardAPI.NewDeleteHandler(
			token.DecodeAuth,
			token.DecodeState,
			teamtable.NewBoardDeleter(db),
			log,
		),
	}))

	taskTitleValidator := taskAPI.NewTitleValidator()
	mux.Handle("/task", api.NewHandler(map[string]api.MethodHandler{
		http.MethodPost: taskAPI.NewPostHandler(
			token.DecodeAuth,
			token.DecodeState,
			taskTitleValidator,
			taskTitleValidator,
			taskAPI.NewColNoValidator(),
			tasktable.NewInserter(db),
			token.EncodeState,
			log,
		),
		http.MethodPatch: taskAPI.NewPatchHandler(
			token.DecodeAuth,
			token.DecodeState,
			taskTitleValidator,
			taskTitleValidator,
			tasktable.NewUpdater(db),
			log,
		),
		http.MethodDelete: taskAPI.NewDeleteHandler(
			token.DecodeAuth,
			token.DecodeState,
			tasktable.NewDeleter(db),
			token.EncodeState,
			log,
		),
	},
	))

	mux.Handle("/tasks", api.NewHandler(map[string]api.MethodHandler{
		http.MethodPatch: tasksAPI.NewPatchHandler(
			token.DecodeAuth,
			token.DecodeState,
			tasksAPI.NewColNoValidator(),
			tasktable.NewMultiUpdater(db),
			token.EncodeState,
			log,
		),
		http.MethodGet: tasksAPI.NewGetHandler(
			token.DecodeAuth,
			tasktable.NewMultiRetriever(db),
			log,
		),
	},
	))

	// serve the registered routes
	log.Info("running server at port " + env.Port)
	if err := http.ListenAndServe(":"+env.Port, mux); err != nil {
		log.Fatal(err.Error())
		os.Exit(5)
	}
}
