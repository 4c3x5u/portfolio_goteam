package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
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
		log.Error(envPort, errPostfix)
	case awsRegion:
		log.Error(envAWSRegion, errPostfix)
	case jwtKey:
		log.Error(envJWTKey, errPostfix)
	case clientOrigin:
		log.Error(envClientOrigin, errPostfix)
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

	mux.Handle("/team", api.NewHandler(map[string]api.MethodHandler{
		http.MethodGet: teamapi.NewGetHandler(
			authDecoder,
			teamtbl.NewRetriever(db),
			teamtbl.NewInserter(db),
			log,
		),
	}))

	mux.Handle("/board", api.NewHandler(map[string]api.MethodHandler{
		http.MethodPost: boardapi.NewPostHandler(
			authDecoder,
			stateDecoder,
			boardapi.NewNameValidator(),
			teamtbl.NewBoardInserter(db),
			stateEncoder,
			log,
		),
		http.MethodPatch: boardapi.NewPatchHandler(
			authDecoder,
			stateDecoder,
			boardapi.NewIDValidator(),
			boardapi.NewNameValidator(),
			teamtbl.NewBoardUpdater(db),
			log,
		),
		http.MethodDelete: boardapi.NewDeleteHandler(
			authDecoder,
			stateDecoder,
			teamtbl.NewBoardDeleter(db),
			stateEncoder,
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
