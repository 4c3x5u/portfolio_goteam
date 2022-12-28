package login

import "net/http"

// Handler is the HTTP handler for the login route.
type Handler struct{}

// NewHandler is the constructor for Handler.
func NewHandler() *Handler {
	return &Handler{}
}

// ServeHTTP responds to requests made to the login route.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
