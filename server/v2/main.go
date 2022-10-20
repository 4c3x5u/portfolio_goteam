package main

import (
	"log"
	"net/http"

	"github.com/kxplxn/goteam/server/v2/api"
)

func main() {
	if err := runWebAPI(); err != nil {
		log.Fatal(err)
	}
}

func runWebAPI() error {
	return serveRoutes(map[string]http.HandlerFunc{
		"/":         api.HandleRoot,
		"/register": api.HandleRegister,
	}, ":1337")
}

func serveRoutes(routes map[string]http.HandlerFunc, port string) error {
	mux := http.NewServeMux()
	for route, handleFunc := range routes {
		mux.HandleFunc(route, handleFunc)
	}
	return http.ListenAndServe(port, mux)
}
