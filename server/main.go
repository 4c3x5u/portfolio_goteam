package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"server/api/login"
	"server/api/register"
	"server/auth"
	"server/db"

	"github.com/golang-jwt/jwt/v4"
)

func main() {
	// Create dependencies that are shared by multiple handlers.
	conn, err := sql.Open("postgres", os.Getenv("DBCONNSTR"))
	if err != nil {
		log.Fatal(err)
	}

	generatorToken := auth.NewGeneratorToken(
		os.Getenv("JWTSIGNATURE"),
		jwt.SigningMethodHS256,
	)

	// Register handlers for API endpoints.
	mux := http.NewServeMux()

	mux.Handle("/register", register.NewHandler(
		register.NewValidatorReq(
			register.NewValidatorUsername(),
			register.NewValidatorPassword(),
		),
		db.NewReaderUser(conn),
		register.NewHasherPwd(),
		db.NewCreatorUser(conn),
		generatorToken,
	))

	mux.Handle("/login", login.NewHandler(
		db.NewReaderUser(conn),
		login.NewComparerHash(),
		generatorToken,
	))

	// Serve the API routes using the ServeMux.
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
