//go:build itest

// Package itest contains integration tests for the application. Each Go file
// except this one and main_test.go corresponds to a HTTP Handler used by the
// application to serve an API endpoint.
package itest

import (
	"database/sql"
	"net/http"
)

// db is the database connection pool used during integration testing.
// It is set in main_test.go/TestMain.
var db *sql.DB

const (
	// jwtKey is the JWT key used for signing and validating JWTs during
	// integration testing.
	jwtKey = "itest-jwt-key-0123456789qwerty"

	// jwtBob123 is the JWT belonging to the bob123 test user. TODO: remove
	jwtBob123 = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ" +
		"ib2IxMjMifQ.Y8_6K50EHUEJlJf4X21fNCFhYWhVIqN3Tw1niz8XwZc"

	// jwtBob124 is the JWT belonging to the bob124 test user. TODO: remove
	jwtBob124 = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ" +
		"ib2IxMjQifQ.LqENrj9APUHgQ3X0HRN6-IFMIg6nyo0_n74KfoxA0qI"

	// JWTs to be used for testing purposes
	jwtTeam1Admin = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZWFtMUF" +
		"kbWluIn0.hdiH2HHc8QFT9VbkpfXKubtV5-mMIT__tmMmYZHMVeA"
	jwtTeam1Member = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZWFtMU" +
		"1lbWJlciJ9.uJbS6vSFZzH1Nfbbto3ega9COg9dMuo63iYHmMYJ6bc"
	jwtTeam2Admin = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZWFtMkF" +
		"kbWluIn0.vjQ93bx9-LK7SZEmhuzISf-Mcf_-A2bZ6VbLn27THPY"
	jwtTeam3Admin = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZWFtM0F" +
		"kbWluIn0.QHFI2okGYug7GNwMwwpwYyTtZkx53I-R-uNjlodCwTU"
)

// addBearerAuth is used in various test cases to authenticate the request
// being sent to a handler.
func addBearerAuth(token string) func(*http.Request) {
	return func(req *http.Request) {
		req.Header.Add("Authorization", "Bearer "+token)
	}
}
