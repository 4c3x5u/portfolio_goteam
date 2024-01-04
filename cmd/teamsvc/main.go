package main

import (
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/joho/godotenv"

	"github.com/kxplxn/goteam/internal/teamsvc/boardapi"
	"github.com/kxplxn/goteam/internal/teamsvc/teamapi"
	"github.com/kxplxn/goteam/pkg/api"
	"github.com/kxplxn/goteam/pkg/cookie"
	"github.com/kxplxn/goteam/pkg/db/teamtbl"
	"github.com/kxplxn/goteam/pkg/log"
)

const (
	// envPort is the name of the environment variable used for setting the port
	// to run the team service on.
	envPort = "TEAM_SERVICE_PORT"

	// envAWSEndpoint is the name of the environment variable used for setting
	// the AWS endpoint to connect to for DynamoDB. It should only be non-empty
	// on local pointing to the local DynamoDB instance.
	envAWSEndpoint = "AWS_ENDPOINT"

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
		awsEndpoint  = os.Getenv(envAWSEndpoint)
		awsAccessKey = os.Getenv(envAWSAccessKey)
		awsSecretKey = os.Getenv(envAWSSecretKey)
		awsRegion    = os.Getenv(envAWSRegion)
		jwtKey       = os.Getenv(envJWTKey)
		clientOrigin = os.Getenv(envClientOrigin)
	)

	// check all environment variables were set
	// - except aws endpoint, which is only set on local
	errPostfix := "was empty"
	switch "" {
	case port:
		log.Error(envPort, errPostfix)
		return
	case awsAccessKey:
		log.Fatal(envAWSAccessKey, errPostfix)
		return
	case awsSecretKey:
		log.Fatal(envAWSSecretKey, errPostfix)
		return
	case awsRegion:
		log.Error(envAWSRegion, errPostfix)
		return
	case jwtKey:
		log.Error(envJWTKey, errPostfix)
		return
	case clientOrigin:
		log.Error(envClientOrigin, errPostfix)
		return
	}

	// define aws config
	cfg := aws.Config{
		Region: awsRegion,
		Credentials: credentials.NewStaticCredentialsProvider(
			awsAccessKey, awsSecretKey, "",
		),
	}
	if awsEndpoint != "" {
		cfg.BaseEndpoint = aws.String(awsEndpoint)
	}

	// create DynamoDB client from config
	db := dynamodb.NewFromConfig(cfg)

	// create auth encoder to be used for authenticating user on all routes
	authDecoder := cookie.NewAuthDecoder([]byte(jwtKey))

	// register handlers for HTTP routes
	mux := http.NewServeMux()

	mux.Handle("/team", api.NewHandler(map[string]api.MethodHandler{
		http.MethodGet: teamapi.NewGetHandler(
			authDecoder,
			teamtbl.NewRetriever(db),
			teamtbl.NewInserter(db),
			teamtbl.NewUpdater(db),
			cookie.NewInviteEncoder([]byte(jwtKey), 1*time.Hour),
			log,
		),
	}))

	mux.Handle("/board", api.NewHandler(map[string]api.MethodHandler{
		http.MethodPost: boardapi.NewPostHandler(
			authDecoder,
			boardapi.NewNameValidator(),
			teamtbl.NewBoardInserter(db),
			log,
		),
		http.MethodPatch: boardapi.NewPatchHandler(
			authDecoder,
			boardapi.NewIDValidator(),
			boardapi.NewNameValidator(),
			teamtbl.NewBoardUpdater(db),
			log,
		),
		http.MethodDelete: boardapi.NewDeleteHandler(
			authDecoder,
			teamtbl.NewBoardDeleter(db),
			log,
		),
	}))

	// serve the registered routes
	log.Info("running team service on port", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
		return
	}
}
