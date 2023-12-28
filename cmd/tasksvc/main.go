package main

import (
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
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

	// envPort is the name of the environment variable used for providing AWS
	// access key to the DynamoDB client.
	envAWSAccessKey = "AWS_ACCESS_KEY"

	// envPort is the name of the environment variable used for providing AWS
	// secret key to the DynamoDB client.
	envAWSSecretKey = "AWS_SECRET_KEY"

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
		awsAccessKey = os.Getenv(envAWSAccessKey)
		awsSecretKey = os.Getenv(envAWSSecretKey)
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
	case awsAccessKey:
		log.Fatal(envAWSAccessKey, errPostfix)
		return
	case awsSecretKey:
		log.Fatal(envAWSSecretKey, errPostfix)
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
	db := dynamodb.NewFromConfig(aws.Config{
		Region: awsRegion,
		Credentials: credentials.NewStaticCredentialsProvider(
			awsAccessKey, awsSecretKey, "",
		),
	})

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
			taskTitleValidator,
			taskTitleValidator,
			tasktbl.NewUpdater(db),
			log,
		),
		http.MethodDelete: taskapi.NewDeleteHandler(
			authDecoder,
			tasktbl.NewDeleter(db),
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
			tasktbl.NewRetrieverByBoard(db),
			authDecoder,
			tasktbl.NewRetrieverByTeam(db),
			stateEncoder,
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
