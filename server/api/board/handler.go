package board

import (
	"net/http"
	"server/auth"
)

// Handler is the http.Handler for the boards route.
type Handler struct{}

// NewHandler creates and returns a new Handler.
func NewHandler() Handler { return Handler{} }

// ServeHTTP responds to requests made to the board route.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Accept only POST requests.
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if _, err := r.Cookie(auth.CookieName); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}
