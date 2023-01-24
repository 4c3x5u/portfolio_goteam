package board

import (
	"net/http"
)

// DELETEHandler can be used to handle the DELETE requests sent to the board
// endpoint.
type DELETEHandler struct{}

// NewDELETEHandler creates and returns a new DELETEHandler.
func NewDELETEHandler() DELETEHandler { return DELETEHandler{} }

// Handle handles the DELETE requests sent to the board endpoint.
func (h DELETEHandler) Handle(
	w http.ResponseWriter, r *http.Request, sub string,
) {
	w.WriteHeader(http.StatusNotImplemented)
	return
}
