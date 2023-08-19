package main

import (
	"errors"
	"os"
)

const (
	// port is the name of the environment variable used for deciding what port
	// to run the server on.
	port = "PORT"

	// dbConnStr is the name of the environment variable used for connecting to
	// the database.
	dbConnStr = "DBCONNSTR"

	// jwtKey is the name of the environment variable used for signing JWTs.
	jwtKey = "JWTKEY"

	// clientOrigin is the name of the environment variable used to set up CORS
	// with the client app.
	clientOrigin = "CLIENTORIGIN"

	// errPostfix is a postfix to environment variable names used for errors
	// returned from the env.validate() function
	errPostfix = " environment variable was empty"
)

// env is used to load, validate and access environment variables.
type env struct {
	Port         string
	DBConnStr    string
	JWTKey       string
	ClientOrigin string
}

// newEnv creates and returns the pointer to a new env, fields of which is
// populated from the currently set environment variables.
func newEnv() *env {
	return &env{
		Port:         os.Getenv(port),
		DBConnStr:    os.Getenv(dbConnStr),
		JWTKey:       os.Getenv(jwtKey),
		ClientOrigin: os.Getenv(clientOrigin),
	}
}

// validate validates that all required environment variables are non-empty.
func (e *env) validate() error {
	switch "" {
	case e.Port:
		return errors.New(port + errPostfix)
	case e.DBConnStr:
		return errors.New(dbConnStr + errPostfix)
	case e.JWTKey:
		return errors.New(jwtKey + errPostfix)
	case e.ClientOrigin:
		return errors.New(clientOrigin + errPostfix)
	default:
		return nil
	}
}
