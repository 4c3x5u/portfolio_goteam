package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	apiRegister "github.com/kxplxn/goteam/server-v2/api/register"
)

func main() {
	db, err := sql.Open("postgres", os.Getenv("DBCONNSTR"))
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.Handle("/register", apiRegister.NewHandler(
		apiRegister.NewCreatorDBUser(db),
		apiRegister.NewValidator(
			apiRegister.NewValidatorUsername(),
			apiRegister.NewValidatorUsername(),
		),
	))

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
