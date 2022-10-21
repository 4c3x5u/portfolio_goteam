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
	return serveRoutes(map[string]http.Handler{
		"/":         &api.HandlerRoot{},
		"/register": &api.HandlerRegister{},
	}, ":1337")
}

func serveRoutes(routes map[string]http.Handler, port string) error {
	mux := http.NewServeMux()
	for route, handler := range routes {
		mux.Handle(route, handler)
	}
	return http.ListenAndServe(port, mux)
}
