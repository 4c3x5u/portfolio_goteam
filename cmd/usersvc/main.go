package main

import (
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
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
		port         = os.Getenv(envSvcPort)
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
		log.Error(envSvcPort, errPostfix)
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
		inviteDecoder = cookie.NewInviteDecoder(key)
		authEncoder   = cookie.NewAuthEncoder(key, dur)
	)

	// register handlers for HTTP routes
	mux := http.NewServeMux()

	mux.Handle("/register", api.NewHandler(map[string]api.MethodHandler{
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

	mux.Handle("/login", api.NewHandler(map[string]api.MethodHandler{
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
		return
	}
}
