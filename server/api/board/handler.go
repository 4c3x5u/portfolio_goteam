package board

import (
	"net/http"

	"server/api"
	"server/auth"
)

// Handler is the http.Handler for the boards route.
type Handler struct {
	tokenValidator auth.TokenValidator
	postHandler    api.MethodHandler
}

// NewHandler creates and returns a new Handler.
func NewHandler(
	tokenValidator auth.TokenValidator,
	postHandler api.MethodHandler,
) Handler {
	return Handler{
		tokenValidator: tokenValidator,
		postHandler:    postHandler,
	}
}

// ServeHTTP responds to requests made to the board route.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only accept the HTTP methods that are handled.
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Get the authentication cookie.
	authCookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.Header().Set(auth.WWWAuthenticate())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Validate authentication cookie value and get sub.
	sub, err := h.tokenValidator.Validate(authCookie.Value)
	if err != nil {
		w.Header().Set(auth.WWWAuthenticate())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Call a MethodHandler based on the HTTP method.
	switch r.Method {
	case http.MethodPost:
		h.postHandler.Handle(w, r, sub)
	}
}
