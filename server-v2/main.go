package main

import (
	"log"
	"net/http"

	apiRegister "github.com/kxplxn/goteam/server-v2/api/register"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/register", apiRegister.NewHandler())

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
