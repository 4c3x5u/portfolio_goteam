package main

import (
	"log"
	"net/http"

	"github.com/kxplxn/goteam/server/v2/api"
	"github.com/kxplxn/goteam/server/v2/relay"
)

func main() {
	if err := runWebAPI(); err != nil {
		log.Fatal(err)
	}
}

func runWebAPI() error {
	apiLogger := relay.NewAPILogger()

	return serveRoutes(map[string]http.Handler{
		"/":         api.NewHandlerRoot(apiLogger),
		"/register": api.NewHandlerRegister(apiLogger),
	}, ":1337")
}

func serveRoutes(routes map[string]http.Handler, port string) error {
	mux := http.NewServeMux()
	for route, handler := range routes {
		mux.Handle(route, handler)
	}
	return http.ListenAndServe(port, mux)
}
