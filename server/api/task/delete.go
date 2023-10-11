package task

import (
	"encoding/json"
	"errors"
	"net/http"
	"server/api"
	pkgLog "server/log"
)

// DELETEHandler is an api.MethodHandler that can be used to handle DELETE
// requests sent to the task route.
type DELETEHandler struct {
	idValidator api.StringValidator
	log         pkgLog.Errorer
}

// NewDELETEHandler creates and returns a new DELETEHandler.
func NewDELETEHandler(
	idValidator api.StringValidator,
	log pkgLog.Errorer,
) DELETEHandler {
	return DELETEHandler{idValidator: idValidator, log: log}
}

// Handle handles the DELETE requests sent to the task route.
func (h DELETEHandler) Handle(
	w http.ResponseWriter, r *http.Request, _ string,
) {
	id := r.URL.Query().Get("id")
	if err := h.idValidator.Validate(id); errors.Is(err, api.ErrStrEmpty) {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(ResBody{
			Error: "Task ID cannot be empty.",
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
			return
		}
	}
	w.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(w).Encode(ResBody{
		Error: "Task ID must be an integer.",
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}
}
