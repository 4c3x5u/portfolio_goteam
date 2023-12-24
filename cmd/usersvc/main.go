package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/joho/godotenv"

	"github.com/kxplxn/goteam/internal/usersvc/loginapi"
	"github.com/kxplxn/goteam/internal/usersvc/registerapi"
	"github.com/kxplxn/goteam/pkg/api"
	"github.com/kxplxn/goteam/pkg/cookie"
	"github.com/kxplxn/goteam/pkg/db/usertbl"
	"github.com/kxplxn/goteam/pkg/log"
)

const (
	// envSvcPort is the name of the environment variable used for setting the
	// port to run the user service on.
	envSvcPort = "USER_SERVICE_PORT"

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
		port         = os.Getenv(envSvcPort)
		awsRegion    = os.Getenv(envAWSRegion)
		jwtKey       = os.Getenv(envJWTKey)
		clientOrigin = os.Getenv(envClientOrigin)
	)

	// check all environment variables were set
	errPostfix := " was empty"
	switch "" {
	case port:
		log.Error(envSvcPort + errPostfix)
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
		inviteDecoder = cookie.NewInviteDecoder(key)
		authEncoder   = cookie.NewAuthEncoder(key, dur)
	)

	// register handlers for HTTP routes
	mux := http.NewServeMux()

	mux.Handle("/user/register", api.NewHandler(map[string]api.MethodHandler{
		http.MethodPost: registerapi.NewPostHandler(
			registerapi.NewUserValidator(
				registerapi.NewUsernameValidator(),
				registerapi.NewPasswordValidator(),
			),
			inviteDecoder,
			registerapi.NewPasswordHasher(),
			usertbl.NewInserter(db),
			authEncoder,
			log,
		),
	}))

	mux.Handle("/user/login", api.NewHandler(map[string]api.MethodHandler{
		http.MethodPost: loginapi.NewPostHandler(
			loginapi.NewValidator(),
			usertbl.NewRetriever(db),
			loginapi.NewPasswordComparator(),
			authEncoder,
			log,
		),
	}))

	// serve the registered routes
	log.Info("running user service on port", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
		os.Exit(5)
	}
}
