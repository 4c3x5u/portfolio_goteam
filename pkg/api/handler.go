package api

import (
	"net/http"
	"os"
)

// MethodHandler describes a type that can be used to serve a certain part of an
// API route that corresponds to a specific HTTP method. It is intended for its
// Handle method to be called after authentication, and with the authenticated
// user's username as the third argument.
type MethodHandler interface {
	Handle(w http.ResponseWriter, r *http.Request, username string)
}

// AllowedMethods takes in a slice of allowed HTTP methods and returns the key
// and the value for the Access-Control-Allow-Methods header.
func AllowedMethods(methods []string) (string, string) {
	if len(methods) == 0 {
		return "", ""
	}

	allowedMethods := ""
	for _, m := range methods {
		allowedMethods += m + ", "
	}

	return "Access-Control-Allow-Methods",
		allowedMethods[:len(allowedMethods)-2]
}

// Handler is a http.Handler that can be used to handle requests.
type Handler struct{ methodHandlers map[string]MethodHandler }

// NewHandler creates and returns a new Handler.
func NewHandler(methodHandlers map[string]MethodHandler) Handler {
	return Handler{methodHandlers: methodHandlers}
}

// ServeHTTP responds to HTTP requests.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// add CORS headers
	w.Header().Set("Access-Control-Allow-Origin", os.Getenv("CLIENTORIGIN"))
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	allowedMethods := []string{http.MethodOptions}
	for method := range h.methodHandlers {
		allowedMethods = append(allowedMethods, method)
	}
	w.Header().Add(AllowedMethods(allowedMethods))

	// if method is OPTIONS, return now with set headers
	if r.Method == http.MethodOptions {
		return
	}

	// find the method handler for the HTTP method of the received request
	for method, methodHandler := range h.methodHandlers {
		if r.Method == method {
			var username string

			// Token sub is used as the username in methodHandler.Handle.
			methodHandler.Handle(w, r, username)
			return
		}
	}

	// if no method handler was found, respond with 405
	w.WriteHeader(http.StatusMethodNotAllowed)
}
