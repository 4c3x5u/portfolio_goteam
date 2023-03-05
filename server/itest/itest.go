//go:build itest

// Package itest contains integration tests for the application. Each Go file
// except this one and main_test.go corresponds to a HTTP Handler used by the
// application to serve an API endpoint.
package itest

import "database/sql"

// dbConn is the database connection pool used during integration testing.
// It is set in main_test.go/TestMain.
var dbConn *sql.DB

// jwtKey is the JWT key used for signing and validating JWTs during integration
// testing.
const jwtKey = "itest-jwt-key-0123456789qwerty"
