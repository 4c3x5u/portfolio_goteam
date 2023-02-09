package board

import (
	"net/http"

	"server/api"
	"server/auth"
)

// Handler is the http.Handler for the boards route.
type Handler struct {
	authHeaderReader   auth.HeaderReader
	authTokenValidator auth.TokenValidator
	methodHandlers     map[string]api.MethodHandler
}

// NewHandler creates and returns a new Handler.
func NewHandler(
	authHeaderReader auth.HeaderReader,
	authTokenValidator auth.TokenValidator,
	methodHandlers map[string]api.MethodHandler,
) Handler {
	return Handler{
		authHeaderReader:   authHeaderReader,
		authTokenValidator: authTokenValidator,
		methodHandlers:     methodHandlers,
	}
}

// ServeHTTP responds to requests made to the board route.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Find the MethodHandler for the HTTP method of the received request.
	for method, methodHandler := range h.methodHandlers {
		// If found, authenticate and handle with MethodHandler.
		if method == r.Method {
			// Get auth token from Authorization header, validate it, and get
			// the subject of the token.
			authToken := h.authHeaderReader.Read(r.Header.Get("Authorization"))
			sub := h.authTokenValidator.Validate(authToken)
			if sub == "" {
				w.Header().Set(auth.WWWAuthenticate())
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// Token sub is used as the username in Handle and this function is
			// returned from.
			methodHandler.Handle(w, r, sub)
			return
		}
	}
	// This path of execution means no MethodHandler was found in
	// h.methodHandlers for the HTTP method of the request.
	w.Header().Add(api.AllowedMethods(http.MethodPost))
	w.WriteHeader(http.StatusMethodNotAllowed)
}
