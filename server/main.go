package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	apiRegister "server/api/register"
)

func main() {
	// todo: use a secret for DBCONNSTR, not an environment variable
	db, err := sql.Open("postgres", os.Getenv("DBCONNSTR"))
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.Handle("/register", apiRegister.NewHandler(
		apiRegister.NewCreatorDBUser(db),
		apiRegister.NewValidator(
			apiRegister.NewValidatorUsername(),
			apiRegister.NewValidatorPassword(),
		),
	))

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
