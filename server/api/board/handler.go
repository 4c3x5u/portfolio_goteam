package board

import (
	"net/http"

	"server/auth"
)

// Handler is the http.Handler for the boards route.
type Handler struct{ authValidator auth.Validator }

// NewHandler creates and returns a new Handler.
func NewHandler(authValidator auth.Validator) Handler {
	return Handler{authValidator: authValidator}
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
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if _, err = h.authValidator.Validate(authCookie.Value); err != nil {
		// This isn't ideal as it will return 401 on both server errors and
		// invaldi JWT. However, due to the way that jwt-go works, this seems
		// to be very hard (impossible?) to avoid.
		// TODO: Figure out something else or use a different jwt library.
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

}
