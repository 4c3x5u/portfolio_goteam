package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	apiRegister "server/api/register"
	"server/db"
)

func main() {
	// todo: use a secret for DBCONNSTR, not an environment variable
	conn, err := sql.Open("postgres", os.Getenv("DBCONNSTR"))
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.Handle("/register", apiRegister.NewHandler(
		apiRegister.NewValidator(
			apiRegister.NewValidatorUsername(),
			apiRegister.NewValidatorPassword(),
		),
		db.NewExistorUser(conn),
		apiRegister.NewHasherPwd(),
		db.NewCreatorUser(conn),
	))

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
