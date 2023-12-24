package main

import (
	"context"
	"net/http"
	"os"
	"time"

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
	"github.com/kxplxn/goteam/pkg/cookie"
	"github.com/kxplxn/goteam/pkg/db/tasktable"
	"github.com/kxplxn/goteam/pkg/db/teamtable"
	"github.com/kxplxn/goteam/pkg/db/usertable"
	"github.com/kxplxn/goteam/pkg/log"
)

const (
	// envServePort is the name of the environment variable used for setting the port
	// to run the server on.
	envServePort = "SERVE_PORT"

	// envAWSRegion is the name of the environment variable used for determining
	// the AWS region to connect to for DynamoDB.
	envAWSRegion = "AWS_REGION"

	// envJWTKey is the name of the environment variable used for signing JWTs.
	envJWTKey = "JWT_KEY"

	// envClientOrigin is the name of the environment variable used to set up CORS
	// with the client app.
	envClientOrigin = "CLIENT_ORIGIN"
)

func main() {
	// create a logger
	log := log.New()

	// load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// get environment variables
	var (
		port         = os.Getenv(envServePort)
		awsRegion    = os.Getenv(envAWSRegion)
		jwtKey       = os.Getenv(envJWTKey)
		clientOrigin = os.Getenv(envClientOrigin)
	)

	// check all environment variables were set
	errPostfix := " was empty"
	switch "" {
	case port:
		log.Error(envServePort + errPostfix)
	case awsRegion:
		log.Error(envAWSRegion + errPostfix)
	case jwtKey:
		log.Error(envJWTKey + errPostfix)
	case clientOrigin:
		log.Error(envClientOrigin + errPostfix)
	}

	// create DynamoDB client
	dbCfg, err := config.LoadDefaultConfig(
		context.Background(), config.WithRegion(awsRegion),
	)
	if err != nil {
		log.Fatal(err)
	}
	db := dynamodb.NewFromConfig(dbCfg)

	// create JWT encoders and decoders
	key := []byte(jwtKey)
	dur := 1 * time.Hour
	var (
		authEncoder   = cookie.NewAuthEncoder(key, dur)
		authDecoder   = cookie.NewAuthDecoder(key)
		stateEncoder  = cookie.NewStateEncoder(key, dur)
		stateDecoder  = cookie.NewStateDecoder(key)
		inviteDecoder = cookie.NewInviteDecoder(key)
	)

	// register handlers for HTTP routes
	mux := http.NewServeMux()

	mux.Handle("/user/register", api.NewHandler(map[string]api.MethodHandler{
		http.MethodPost: registerAPI.NewPostHandler(
			registerAPI.NewUserValidator(
				registerAPI.NewUsernameValidator(),
				registerAPI.NewPasswordValidator(),
			),
			inviteDecoder,
			registerAPI.NewPasswordHasher(),
			usertable.NewInserter(db),
			authEncoder,
			log,
		),
	}))

	mux.Handle("/user/login", api.NewHandler(map[string]api.MethodHandler{
		http.MethodPost: loginAPI.NewPostHandler(
			loginAPI.NewValidator(),
			usertable.NewRetriever(db),
			loginAPI.NewPasswordComparator(),
			authEncoder,
			log,
		),
	}))

	mux.Handle("/team", api.NewHandler(map[string]api.MethodHandler{
		http.MethodGet: teamAPI.NewGetHandler(
			authDecoder,
			teamtable.NewRetriever(db),
			teamtable.NewInserter(db),
			log,
		),
	}))

	mux.Handle("/team/board", api.NewHandler(map[string]api.MethodHandler{
		http.MethodPost: boardAPI.NewPostHandler(
			authDecoder,
			stateDecoder,
			boardAPI.NewNameValidator(),
			teamtable.NewBoardInserter(db),
			stateEncoder,
			log,
		),
		http.MethodPatch: boardAPI.NewPatchHandler(
			authDecoder,
			stateDecoder,
			boardAPI.NewIDValidator(),
			boardAPI.NewNameValidator(),
			teamtable.NewBoardUpdater(db),
			log,
		),
		http.MethodDelete: boardAPI.NewDeleteHandler(
			authDecoder,
			stateDecoder,
			teamtable.NewBoardDeleter(db),
			stateEncoder,
			log,
		),
	}))

	taskTitleValidator := taskAPI.NewTitleValidator()
	mux.Handle("/task", api.NewHandler(map[string]api.MethodHandler{
		http.MethodPost: taskAPI.NewPostHandler(
			authDecoder,
			stateDecoder,
			taskTitleValidator,
			taskTitleValidator,
			taskAPI.NewColNoValidator(),
			tasktable.NewInserter(db),
			stateEncoder,
			log,
		),
		http.MethodPatch: taskAPI.NewPatchHandler(
			authDecoder,
			stateDecoder,
			taskTitleValidator,
			taskTitleValidator,
			tasktable.NewUpdater(db),
			log,
		),
		http.MethodDelete: taskAPI.NewDeleteHandler(
			authDecoder,
			stateDecoder,
			tasktable.NewDeleter(db),
			stateEncoder,
			log,
		),
	},
	))

	mux.Handle("/tasks", api.NewHandler(map[string]api.MethodHandler{
		http.MethodPatch: tasksAPI.NewPatchHandler(
			authDecoder,
			stateDecoder,
			tasksAPI.NewColNoValidator(),
			tasktable.NewMultiUpdater(db),
			stateEncoder,
			log,
		),
		http.MethodGet: tasksAPI.NewGetHandler(
			authDecoder,
			tasktable.NewMultiRetriever(db),
			log,
		),
	},
	))

	// serve the registered routes
	log.Info("running server at port " + port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
		os.Exit(5)
	}
}
