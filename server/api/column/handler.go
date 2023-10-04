package column

import (
	"net/http"
	"server/api"
)

// Handler is a http.Handler that can be used to handle column requests.
type Handler struct{}

// NewHandler creates and returns a new Handler.
func NewHandler() Handler { return Handler{} }

// ServeHTTP responds to requests made to the column route.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only allow PATCH requests.
	if r.Method != http.MethodPatch {
		w.Header().Add(api.AllowedMethods(http.MethodPost))
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	w.WriteHeader(http.StatusOK)
	return
}
