package task

import (
	"encoding/json"
	"net/http"
	"server/api"
)

// ReqBody defines the request body for requests handled by method handlers.
type ReqBody struct {
	Title string `json:"title"`
}

// ResBody defines the response body for requests handled by method handlers.
type ResBody struct {
	Error string `json:"error"`
}

// POSTHandler is an api.MethodHandler that can be used to handle POST task
// requests.
type POSTHandler struct {
	titleValidator api.StringValidator
}

// NewPOSTHandler creates and returns a new POSTHandler.
func NewPOSTHandler(titleValidator api.StringValidator) *POSTHandler {
	return &POSTHandler{titleValidator: titleValidator}
}

// Handle handles the POST requests sent to the task route.
func (h *POSTHandler) Handle(
	w http.ResponseWriter, r *http.Request, _ string,
) {
	var reqBody ReqBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := h.titleValidator.Validate(reqBody.Title); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if encodeErr := json.NewEncoder(w).Encode(ResBody{
			Error: err.Error(),
		}); encodeErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
}
