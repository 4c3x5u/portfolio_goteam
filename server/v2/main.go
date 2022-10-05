package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		devRelay(w,
			"health check status: OK\n"+
				"available endpoints:\n"+
				"/register",
		)
	})

	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		devRelay(w, "register endpoint is hit")
	})

	if err := http.ListenAndServe(":1337", mux); err != nil {
		log.Fatal(err)
	}
}

func devRelay(w http.ResponseWriter, msg string) {
	log.Println(msg)
	if _, err := w.Write([]byte(msg)); err != nil {
		log.Fatal(err)
	}
}
