package subtask

import (
	"encoding/json"
	"net/http"
)

// ResBody defines the response body for requests handled by PATCHHandler.
type ResBody struct {
	Error string `json:"error"`
}

// PATCHHandler is an api.MethodHandler that can be used to handle PATCH
// requests sent to the subtask route.
type PATCHHandler struct{}

// NewPATCHHandler creates and returns a new PATCHandler.
func NewPATCHHandler() PATCHHandler { return PATCHHandler{} }

// Handle handles the PATCH requests sent to the subtask route.
func (h PATCHHandler) Handle(w http.ResponseWriter, r *http.Request, _ string) {
	w.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(w).Encode(
		ResBody{Error: "Subtask ID cannot be empty."},
	); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
