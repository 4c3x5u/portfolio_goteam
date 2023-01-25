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
	isHandled := false
	for method, methodHandler := range h.methodHandlers {
		if method == r.Method {
			// Get auth token from Authorization header, validate it, and get the
			// subject of the token.
			authToken := h.authHeaderReader.Read(r.Header.Get("Authorization"))
			sub := h.authTokenValidator.Validate(authToken)
			if sub == "" {
				w.Header().Set(auth.WWWAuthenticate())
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			isHandled = true

			methodHandler.Handle(w, r, sub)
		}
	}
	if !isHandled {
		w.Header().Add(api.AllowedMethods(http.MethodPost))
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
