package api

import (
	"net/http"

	"server/auth"
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
type Handler struct {
	authHeaderReader   auth.HeaderReader
	authTokenValidator auth.TokenValidator
	methodHandlers     map[string]MethodHandler
}

// NewHandler creates and returns a new Handler.
func NewHandler(
	authHeaderReader auth.HeaderReader,
	authTokenValidator auth.TokenValidator,
	methodHandlers map[string]MethodHandler,
) Handler {
	return Handler{
		authHeaderReader:   authHeaderReader,
		authTokenValidator: authTokenValidator,
		methodHandlers:     methodHandlers,
	}
}

// ServeHTTP responds to HTTP requests.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Keep track of allowed methods to return them in
	// "Access-Control-Allow-Methods" header on 405.
	var allowedMethods []string

	// Find the MethodHandler for the HTTP method of the received request.
	for method, methodHandler := range h.methodHandlers {
		allowedMethods = append(allowedMethods, method)

		// If found, authenticate and handle with MethodHandler.
		if r.Method == method {
			// Get auth token from Authorization header, validate it, and get
			// the subject of the token.
			authToken := h.authHeaderReader.Read(
				r.Header.Get(auth.AuthorizationHeader),
			)
			sub := h.authTokenValidator.Validate(authToken)
			if sub == "" {
				w.Header().Set(auth.WWWAuthenticate())
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// Token sub is used as the username in methodHandler.Handle.
			methodHandler.Handle(w, r, sub)
			return
		}
	}
	// This path of execution means no MethodHandler was found in
	// h.methodHandlers for the HTTP method of the request.
	w.Header().Add(AllowedMethods(allowedMethods))
	w.WriteHeader(http.StatusMethodNotAllowed)
}
