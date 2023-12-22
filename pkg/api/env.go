// Package env contains code for reading, accessing, and validating environment
// variables.
// TODO: split into separate files/packages for each internal package.
package api

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

// Env is used to load, validate and access environment variables.
type Env struct {
	Port         string
	DBConnStr    string
	JWTKey       string
	ClientOrigin string
}

// NewEnv creates and returns the pointer to a new env, fields of which is
// populated from the currently set environment variables.
func NewEnv() *Env {
	return &Env{
		Port:         os.Getenv(port),
		DBConnStr:    os.Getenv(dbConnStr),
		JWTKey:       os.Getenv(jwtKey),
		ClientOrigin: os.Getenv(clientOrigin),
	}
}

// Validate validates that all required environment variables are non-empty.
func (e *Env) Validate() error {
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
