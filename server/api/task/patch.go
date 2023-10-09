package task

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"server/api"
	"server/dbaccess"
	taskTable "server/dbaccess/task"
	pkgLog "server/log"
)

// PATCHHandler is an api.MethodHandler that can be used to handle PATCH
// requests sent to the task route.
type PATCHHandler struct {
	idValidator           api.StringValidator
	taskTitleValidator    api.StringValidator
	subtaskTitleValidator api.StringValidator
	taskSelector          dbaccess.Selector[taskTable.Record]
	log                   pkgLog.Errorer
}

// NewPATCHHandler creates and returns a new PATCHHandler.
func NewPATCHHandler(
	idValidator api.StringValidator,
	taskTitleValidator api.StringValidator,
	subtaskTitleValidator api.StringValidator,
	taskSelector dbaccess.Selector[taskTable.Record],
	log pkgLog.Errorer,
) *PATCHHandler {
	return &PATCHHandler{
		idValidator:           idValidator,
		taskTitleValidator:    taskTitleValidator,
		subtaskTitleValidator: subtaskTitleValidator,
		taskSelector:          taskSelector,
		log:                   log,
	}
}

// Handle handles the PATCH requests sent to the task route.
func (h *PATCHHandler) Handle(
	w http.ResponseWriter, r *http.Request, _ string,
) {
	id := r.URL.Query().Get("id")
	if err := h.idValidator.Validate(id); err != nil {
		var errMsg string
		if errors.Is(err, api.ErrStrEmpty) {
			errMsg = "Task ID cannot be empty."
		} else if errors.Is(err, api.ErrStrNotInt) {
			errMsg = "Task ID must be an integer."
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

	var reqBody ReqBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}

	// Validate task title.
	if err := h.taskTitleValidator.Validate(reqBody.Title); err != nil {
		var errMsg string
		if errors.Is(err, api.ErrStrEmpty) {
			errMsg = "Task title cannot be empty."
		} else if errors.Is(err, api.ErrStrTooLong) {
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

	// Validate subtask titles
	for _, title := range reqBody.SubtaskTitles {
		if err := h.subtaskTitleValidator.Validate(title); err != nil {
			var errMsg string
			if errors.Is(err, api.ErrStrEmpty) {
				errMsg = "Subtask title cannot be empty."
			} else if errors.Is(err, api.ErrStrTooLong) {
				errMsg = "Subtask title cannot be longer than 50 characters."
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

	// Find the task in the database to get its columnID.
	_, err := h.taskSelector.Select(id)
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(ResBody{
			Error: "Task not found.",
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}
}
