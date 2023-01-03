package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	apiLogin "server/api/login"
	apiRegister "server/api/register"
	"server/db"
)

func main() {
	conn, err := sql.Open("postgres", os.Getenv("DBCONNSTR"))
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.Handle("/register", apiRegister.NewHandler(
		apiRegister.NewValidatorReq(
			apiRegister.NewValidatorUsername(),
			apiRegister.NewValidatorPassword(),
		),
		db.NewReaderUser(conn),
		apiRegister.NewHasherPwd(),
		db.NewCreatorUser(conn),
		db.NewCreatorSession(conn),
	))

	mux.Handle("/login", apiLogin.NewHandler(
		db.NewReaderUser(conn),
		apiLogin.NewComparerHash(),
		db.NewReaderSession(conn),
	))

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
