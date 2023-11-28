package api

import (
	"errors"
	"net/http"
	"os"

	"github.com/kxplxn/goteam/server/auth"
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
	bearerTokenReader  auth.HeaderReader //TODO: remove
	authTokenValidator auth.TokenValidator
	methodHandlers     map[string]MethodHandler
}

// NewHandler creates and returns a new Handler.
func NewHandler(
	bearerTokenReader auth.HeaderReader, //TODO: remove
	authTokenValidator auth.TokenValidator,
	methodHandlers map[string]MethodHandler,
) Handler {
	return Handler{
		bearerTokenReader:  bearerTokenReader, //TODO: remove
		authTokenValidator: authTokenValidator,
		methodHandlers:     methodHandlers,
	}
}

// ServeHTTP responds to HTTP requests.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Add CORS headers.
	w.Header().Set("Access-Control-Allow-Origin", os.Getenv("CLIENTORIGIN"))
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	allowedMethods := []string{http.MethodOptions}
	for method := range h.methodHandlers {
		allowedMethods = append(allowedMethods, method)
	}
	w.Header().Add(AllowedMethods(allowedMethods))
	// If method OPTIONS, return with set headers.
	if r.Method == http.MethodOptions {
		return
	}

	// Find the MethodHandler for the HTTP method of the received request.
	for method, methodHandler := range h.methodHandlers {
		if r.Method == method {
			var username string

			// If it's NOT a register or login request, validate JWT.
			url := r.URL.String()
			srvOrigin := os.Getenv("SERVERORIGIN")
			if url != srvOrigin+"/register" && url != srvOrigin+"/login" {
				// Get auth token from Authorization header, validate it, and get
				// the subject of the token.
				ck, err := r.Cookie(auth.CookieName)
				if errors.Is(err, http.ErrNoCookie) {
					w.Header().Set(auth.WWWAuthenticate())
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				username = h.authTokenValidator.Validate(ck.Value)
				if username == "" {
					w.Header().Set(auth.WWWAuthenticate())
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
			}

			// Token sub is used as the username in methodHandler.Handle.
			methodHandler.Handle(w, r, username)
			return
		}
	}
	// This path of execution means no MethodHandler was found in
	// h.methodHandlers for the HTTP method of the request.
	w.WriteHeader(http.StatusMethodNotAllowed)
}
