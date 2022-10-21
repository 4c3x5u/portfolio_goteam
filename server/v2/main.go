package main

import (
	"log"
	"net/http"

	"github.com/kxplxn/goteam/server/v2/api/handlers"
	"github.com/kxplxn/goteam/server/v2/relay"
)

func main() {
	if err := runWebAPI(); err != nil {
		log.Fatal(err)
	}
}

func runWebAPI() error {
	apiLogger := relay.NewAPILogger()

	return handleRoutes(map[string]http.Handler{
		"/":         handlers.NewRoot(apiLogger),
		"/register": handlers.NewRegister(apiLogger),
	}, ":1337")
}

func handleRoutes(routeHandlers map[string]http.Handler, port string) error {
	mux := http.NewServeMux()
	for route, handler := range routeHandlers {
		mux.Handle(route, handler)
	}
	return http.ListenAndServe(port, mux)
}
