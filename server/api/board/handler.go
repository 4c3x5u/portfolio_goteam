package board

import "net/http"

// Handler is the http.Handler for the boards route.
type Handler struct{}

func NewHandler() Handler { return Handler{} }

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Accept only POST requests.
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
