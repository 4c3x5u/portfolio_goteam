package task

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"server/api"
	"server/dbaccess"
	columnTable "server/dbaccess/column"
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
	columnSelector        dbaccess.Selector[columnTable.Record]
	userBoardSelector     dbaccess.RelSelector[bool]
	log                   pkgLog.Errorer
}

// NewPATCHHandler creates and returns a new PATCHHandler.
func NewPATCHHandler(
	idValidator api.StringValidator,
	taskTitleValidator api.StringValidator,
	subtaskTitleValidator api.StringValidator,
	taskSelector dbaccess.Selector[taskTable.Record],
	columnSelector dbaccess.Selector[columnTable.Record],
	userBoardSelector dbaccess.RelSelector[bool],
	log pkgLog.Errorer,
) *PATCHHandler {
	return &PATCHHandler{
		idValidator:           idValidator,
		taskTitleValidator:    taskTitleValidator,
		subtaskTitleValidator: subtaskTitleValidator,
		taskSelector:          taskSelector,
		columnSelector:        columnSelector,
		userBoardSelector:     userBoardSelector,
		log:                   log,
	}
}

// Handle handles the PATCH requests sent to the task route.
func (h *PATCHHandler) Handle(
	w http.ResponseWriter, r *http.Request, username string,
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

	// AUTHORIZATION:
	// To authorise this user, we must check that both the source column and the
	// target column belong to a board that the user is the admin of. The ID to
	// the source column can be retrieved from the task table, and the ID to the
	// target column is retrieved within the request body.

	// Select the task in the database to get its columnID.
	task, err := h.taskSelector.Select(id)
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

	for _, columnID := range []int{task.ColumnID, reqBody.ColumnID} {
		// Select the column from the database with the columnID to get the
		// board ID.
		column, err := h.columnSelector.Select(strconv.Itoa(columnID))
		if err != nil {
			// Return 500 on any error (even sql.ErrNoRows) because if task was
			// found, so must the column because the columnID is a foreign key for
			// the column table.
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
			return
		}

		// Select the isAdmin column of the user-board relationship record from
		// the database with the user's username and the column's board ID.
		isAdmin, err := h.userBoardSelector.Select(
			username, strconv.Itoa(column.BoardID),
		)
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusUnauthorized)
			if err := json.NewEncoder(w).Encode(ResBody{
				Error: "You do not have access to this board.",
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
		if !isAdmin {
			w.WriteHeader(http.StatusUnauthorized)
			if err := json.NewEncoder(w).Encode(ResBody{
				Error: "Only board admins can edit tasks.",
			}); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				h.log.Error(err.Error())
			}
			return
		}
	}
}
