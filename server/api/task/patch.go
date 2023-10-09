package task

import (
	"encoding/json"
	"net/http"
	"server/api"
	pkgLog "server/log"
)

// PATCHHandler is an api.MethodHandler that can be used to handle PATCH
// requests sent to the task route.
type PATCHHandler struct {
	taskTitleValidator api.StringValidator
	log                pkgLog.Errorer
}

// NewPATCHHandler creates and returns a new PATCHHandler.
func NewPATCHHandler(
	taskTitleValidator api.StringValidator, log pkgLog.Errorer,
) *PATCHHandler {
	return &PATCHHandler{taskTitleValidator: taskTitleValidator, log: log}
}

// Handle handles the PATCH requests sent to the task route.
func (h *PATCHHandler) Handle(
	w http.ResponseWriter, r *http.Request, _ string,
) {
	w.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(w).Encode(ResBody{
		Error: "Task title cannot be empty.",
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
	}
	return
}
