package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/joho/godotenv"

	"github.com/kxplxn/goteam/internal/tasksvc/taskapi"
	"github.com/kxplxn/goteam/internal/tasksvc/tasksapi"
	"github.com/kxplxn/goteam/pkg/api"
	"github.com/kxplxn/goteam/pkg/cookie"
	"github.com/kxplxn/goteam/pkg/db/tasktbl"
	"github.com/kxplxn/goteam/pkg/log"
)

const (
	// envPort is the name of the environment variable used for setting the port
	// to run the task service on.
	envPort = "TASK_SERVICE_PORT"

	// envAWSRegion is the name of the environment variable used for determining
	// the AWS region to connect to for DynamoDB.
	envAWSRegion = "AWS_REGION"

	// envJWTKey is the name of the environment variable used for signing JWTs.
	envJWTKey = "JWT_KEY"

	// envClientOrigin is the name of the environment variable used to set up
	// CORS with the client app.
	envClientOrigin = "CLIENT_ORIGIN"
)

func main() {
	// create a logger
	log := log.New()

	// load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
		return
	}

	// get environment variables
	var (
		port         = os.Getenv(envPort)
		awsRegion    = os.Getenv(envAWSRegion)
		jwtKey       = os.Getenv(envJWTKey)
		clientOrigin = os.Getenv(envClientOrigin)
	)

	// check all environment variables were set
	errPostfix := "was empty"
	switch "" {
	case port:
		log.Fatal(envPort, errPostfix)
		return
	case awsRegion:
		log.Fatal(envAWSRegion, errPostfix)
		return
	case jwtKey:
		log.Fatal(envJWTKey, errPostfix)
		return
	case clientOrigin:
		log.Fatal(envClientOrigin, errPostfix)
		return
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
		authDecoder  = cookie.NewAuthDecoder(key)
		stateEncoder = cookie.NewStateEncoder(key, dur)
		stateDecoder = cookie.NewStateDecoder(key)
	)

	// register handlers for HTTP routes
	mux := http.NewServeMux()

	taskTitleValidator := taskapi.NewTitleValidator()
	mux.Handle("/task", api.NewHandler(map[string]api.MethodHandler{
		http.MethodPost: taskapi.NewPostHandler(
			authDecoder,
			stateDecoder,
			taskTitleValidator,
			taskTitleValidator,
			taskapi.NewColNoValidator(),
			tasktbl.NewInserter(db),
			stateEncoder,
			log,
		),
		http.MethodPatch: taskapi.NewPatchHandler(
			authDecoder,
			stateDecoder,
			taskTitleValidator,
			taskTitleValidator,
			tasktbl.NewUpdater(db),
			log,
		),
		http.MethodDelete: taskapi.NewDeleteHandler(
			authDecoder,
			stateDecoder,
			tasktbl.NewDeleter(db),
			stateEncoder,
			log,
		),
	}))

	mux.Handle("/tasks", api.NewHandler(map[string]api.MethodHandler{
		http.MethodPatch: tasksapi.NewPatchHandler(
			authDecoder,
			stateDecoder,
			tasksapi.NewColNoValidator(),
			tasktbl.NewMultiUpdater(db),
			stateEncoder,
			log,
		),
		http.MethodGet: tasksapi.NewGetHandler(
			tasksapi.NewBoardIDValidator(),
			stateDecoder,
			tasktbl.NewMultiRetriever(db),
			log,
		),
	}))

	// serve the registered routes
	log.Info("running task service on port", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
		return
	}
}
