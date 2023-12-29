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

// Handler is a http.Handler that can be used to handle requests.
type Handler struct{ methodHandlers map[string]MethodHandler }

// NewHandler creates and returns a new Handler.
func NewHandler(methodHandlers map[string]MethodHandler) Handler {
	return Handler{methodHandlers: methodHandlers}
}

// ServeHTTP responds to HTTP requests.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// add cors headers
	w.Header().Set("Access-Control-Allow-Origin", os.Getenv("CLIENTORIGIN"))
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Add("Access-Control-Allow-Credentials", "true")

	// add allowed methods header
	allowedMethods := make([]string, len(h.methodHandlers)+1)
	allowedMethods[0] = http.MethodOptions
    for method := range h.methodHandlers {
		allowedMethods = append(allowedMethods, method)
	}
	w.Header().Add(allowedMethodsHeader(allowedMethods))

	// if method is OPTIONS, return now with set headers
	if r.Method == http.MethodOptions {
		return
	}

	// get method handler and handle request with it
	methodHandler, ok := h.methodHandlers[r.Method]
	if !ok {
		// return 405 if no method handler corresponds to the given method
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	methodHandler.Handle(w, r, "")
}

// allowedMethodsHeader takes in a slice of allowed HTTP methods and returns the
// key and the value for the Access-Control-Allow-Methods header.
func allowedMethodsHeader(methods []string) (string, string) {
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
