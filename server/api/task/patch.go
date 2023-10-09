package task

import (
	"encoding/json"
	"errors"
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
	var reqBody ReqBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}

	if err := h.taskTitleValidator.Validate(reqBody.Title); err != nil {
		var errMsg string
		if errors.Is(err, errTitleEmpty) {
			errMsg = "Task title cannot be empty."
		} else if errors.Is(err, errTitleTooLong) {
			errMsg = "Task title cannot be longer than 50 characters."
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(ResBody{
			Error: errMsg,
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	}
}
