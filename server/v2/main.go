package main

import (
	"net/http"

	"github.com/kxplxn/goteam/server/v2/api"
	"github.com/kxplxn/goteam/server/v2/log"
)

func main() {
	log.ErrToConsole(runWebAPI())
}

func runWebAPI() error {
	return serveRoutes(map[string]http.HandlerFunc{
		"/":         api.ServeRoot,
		"/register": api.ServeRegister,
	}, ":1337")
}

func serveRoutes(routes map[string]http.HandlerFunc, port string) error {
	mux := http.NewServeMux()
	for route, handleFunc := range routes {
		mux.HandleFunc(route, handleFunc)
	}
	return http.ListenAndServe(port, mux)
}
