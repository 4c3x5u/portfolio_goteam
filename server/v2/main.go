package main

import (
	"log"
	"net/http"
)

func main() {
	if err := runWebAPI(); err != nil {
		log.Fatal(err)
	}
}

func runWebAPI() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		relay := NewLogger(w)
		relay.Msg(
			"health check status: OK\n" +
				"available endpoints:\n" +
				"/register",
		)
	})

	mux.Handle("/register", &RegisterHandler{})

	return http.ListenAndServe(":1337", mux)
}
