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
	conn, err := sql.Open("postgres", os.Getenv("DBCONNSTR"))
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.Handle("/register", register.NewHandler(
		register.NewValidatorReq(
			register.NewValidatorUsername(),
			register.NewValidatorPassword(),
		),
		db.NewReaderUser(conn),
		register.NewHasherPwd(),
		db.NewCreatorUser(conn),
		auth.NewGeneratorToken(os.Getenv("JWTSIGNATURE"), jwt.SigningMethodHS256),
	))

	mux.Handle("/login", login.NewHandler(
		db.NewReaderUser(conn),
		login.NewComparerHash(),
		db.NewUpserterSession(conn),
	))

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
