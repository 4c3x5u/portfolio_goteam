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
	postHandler        api.MethodHandler
	deleteHandler      api.MethodHandler
}

// NewHandler creates and returns a new Handler.
func NewHandler(
	authHeaderReader auth.HeaderReader,
	authTokenValidator auth.TokenValidator,
	postHandler api.MethodHandler,
	deleteHandler api.MethodHandler,
) Handler {
	return Handler{
		authHeaderReader:   authHeaderReader,
		authTokenValidator: authTokenValidator,
		postHandler:        postHandler,
		deleteHandler:      deleteHandler,
	}
}

// ServeHTTP responds to requests made to the board route.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only accept the HTTP methods that are handled.
	if r.Method != http.MethodPost && r.Method != http.MethodDelete {
		w.Header().Add(api.AllowedMethods(http.MethodPost))
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Get auth token from Authorization header, validate it, and get the
	// subject of the token.
	authToken := h.authHeaderReader.Read(r.Header.Get("Authorization"))
	sub := h.authTokenValidator.Validate(authToken)
	if sub == "" {
		w.Header().Set(auth.WWWAuthenticate())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Call a MethodHandler based on the HTTP method.
	switch r.Method {
	case http.MethodPost:
		h.postHandler.Handle(w, r, sub)
	case http.MethodDelete:
		h.deleteHandler.Handle(w, r, sub)
	}
}
