package board

import (
	"net/http"

	"server/auth"
)

// Handler is the http.Handler for the boards route.
type Handler struct{ tokenValidator auth.TokenValidator }

// NewHandler creates and returns a new Handler.
func NewHandler(tokenValidator auth.TokenValidator) Handler {
	return Handler{tokenValidator: tokenValidator}
}

// ServeHTTP responds to requests made to the board route.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Accept only POST requests.
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Get and validate the authentication cookie.
	authCookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.Header().Set(auth.WWWAuthenticate())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if _, err = h.tokenValidator.Validate(authCookie.Value); err != nil {
		w.Header().Set(auth.WWWAuthenticate())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}
