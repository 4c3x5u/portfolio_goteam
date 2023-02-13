//go:build itest

// Package itest contains integration tests for the application.
package itest

// serverURL is the url that is used to send requests to the server running in
// the test container. It is used during setup in main_test.go/TestMain.
const serverURL = "http://localhost:8081"

// dbConnStr is the connection string for the test db.
const dbConnStr = "postgres://itestuser:itestpwd@localhost:5432/itestdb?sslmode=disable"
