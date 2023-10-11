package task

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/kxplxn/goteam/server/api"
	"github.com/kxplxn/goteam/server/dbaccess"
	columnTable "github.com/kxplxn/goteam/server/dbaccess/column"
	taskTable "github.com/kxplxn/goteam/server/dbaccess/task"
	pkgLog "github.com/kxplxn/goteam/server/log"
)

// PATCHReqBody defines the request body for requests handled by PATCHHandler.
type PATCHReqBody struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Subtasks    []struct {
		Title  string `json:"title"`
		Order  int    `json:"order"`
		IsDone bool   `json:"done"`
	} `json:"subtasks"`
}

// PATCHHandler is an api.MethodHandler that can be used to handle PATCH
// requests sent to the task route.
type PATCHHandler struct {
	idValidator           api.StringValidator
	taskTitleValidator    api.StringValidator
	subtaskTitleValidator api.StringValidator
	taskSelector          dbaccess.Selector[taskTable.Record]
	columnSelector        dbaccess.Selector[columnTable.Record]
	userBoardSelector     dbaccess.RelSelector[bool]
	taskUpdater           dbaccess.Updater[taskTable.UpRecord]
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
	taskUpdater dbaccess.Updater[taskTable.UpRecord],
	log pkgLog.Errorer,
) *PATCHHandler {
	return &PATCHHandler{
		idValidator:           idValidator,
		taskTitleValidator:    taskTitleValidator,
		subtaskTitleValidator: subtaskTitleValidator,
		taskSelector:          taskSelector,
		columnSelector:        columnSelector,
		userBoardSelector:     userBoardSelector,
		taskUpdater:           taskUpdater,
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

	var reqBody PATCHReqBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}

	// Validate task title and start building a db-insertable task record.
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

	// Validate subtask titles and transform them into db-insertable types.
	var subtaskRecords []taskTable.Subtask
	for _, subtask := range reqBody.Subtasks {
		if err := h.subtaskTitleValidator.Validate(subtask.Title); err != nil {
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
		subtaskRecords = append(
			subtaskRecords,
			taskTable.NewSubtask(subtask.Title, subtask.Order, subtask.IsDone),
		)
	}

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

	// Get the column from the database to access its board ID for
	// authorization.
	column, err := h.columnSelector.Select(strconv.Itoa(task.ColumnID))
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
		w.WriteHeader(http.StatusForbidden)
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
		w.WriteHeader(http.StatusForbidden)
		if err := json.NewEncoder(w).Encode(ResBody{
			Error: "Only board admins can edit tasks.",
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	}

	// Update the task and subtasks in the database.
	if err := h.taskUpdater.Update(id, taskTable.NewUpRecord(
		reqBody.Title, reqBody.Description, subtaskRecords,
	)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}
}
