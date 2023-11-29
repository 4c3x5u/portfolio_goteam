package task

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/kxplxn/goteam/server/api"
	"github.com/kxplxn/goteam/server/dbaccess"
	boardTable "github.com/kxplxn/goteam/server/dbaccess/board"
	columnTable "github.com/kxplxn/goteam/server/dbaccess/column"
	taskTable "github.com/kxplxn/goteam/server/dbaccess/task"
	userTable "github.com/kxplxn/goteam/server/dbaccess/user"
	pkgLog "github.com/kxplxn/goteam/server/log"
)

// PATCHReq defines the body of PATCH task requests.
type PATCHReq struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Subtasks    []struct {
		Title  string `json:"title"`
		Order  int    `json:"order"`
		IsDone bool   `json:"done"`
	} `json:"subtasks"`
}

// PATCHResp defines the body of PATCH task responses.
type PATCHResp struct {
	Error string `json:"error"`
}

// PATCHHandler is an api.MethodHandler that can be used to handle PATCH
// requests sent to the task route.
type PATCHHandler struct {
	userSelector          dbaccess.Selector[userTable.Record]
	idValidator           api.StringValidator
	taskTitleValidator    api.StringValidator
	subtaskTitleValidator api.StringValidator
	taskSelector          dbaccess.Selector[taskTable.Record]
	columnSelector        dbaccess.Selector[columnTable.Record]
	boardSelector         dbaccess.Selector[boardTable.Record]
	taskUpdater           dbaccess.Updater[taskTable.UpRecord]
	log                   pkgLog.Errorer
}

// NewPATCHHandler creates and returns a new PATCHHandler.
func NewPATCHHandler(
	userSelector dbaccess.Selector[userTable.Record],
	idValidator api.StringValidator,
	taskTitleValidator api.StringValidator,
	subtaskTitleValidator api.StringValidator,
	taskSelector dbaccess.Selector[taskTable.Record],
	columnSelector dbaccess.Selector[columnTable.Record],
	boardSelector dbaccess.Selector[boardTable.Record],
	taskUpdater dbaccess.Updater[taskTable.UpRecord],
	log pkgLog.Errorer,
) *PATCHHandler {
	return &PATCHHandler{
		userSelector:          userSelector,
		idValidator:           idValidator,
		taskTitleValidator:    taskTitleValidator,
		subtaskTitleValidator: subtaskTitleValidator,
		taskSelector:          taskSelector,
		columnSelector:        columnSelector,
		boardSelector:         boardSelector,
		taskUpdater:           taskUpdater,
		log:                   log,
	}
}

// Handle handles the PATCH requests sent to the task route.
func (h *PATCHHandler) Handle(
	w http.ResponseWriter, r *http.Request, username string,
) {
	// Validate that the user is a team admin..
	user, err := h.userSelector.Select(username)
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusUnauthorized)
		if err := json.NewEncoder(w).Encode(PATCHResp{
			Error: "Username is not recognised.",
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
	if !user.IsAdmin {
		w.WriteHeader(http.StatusForbidden)
		if err := json.NewEncoder(w).Encode(PATCHResp{
			Error: "Only team admins can edit tasks.",
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	}

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
		if err := json.NewEncoder(w).Encode(PATCHResp{
			Error: errMsg,
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	}

	var reqBody PATCHReq
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
		if err := json.NewEncoder(w).Encode(PATCHResp{
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
			if err := json.NewEncoder(w).Encode(PATCHResp{
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
		if err := json.NewEncoder(w).Encode(PATCHResp{
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

	// Validate that the board belongs to the team that the user is the admin
	// of.
	board, err := h.boardSelector.Select(strconv.Itoa(column.BoardID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return

	}
	if board.TeamID != user.TeamID {
		w.WriteHeader(http.StatusForbidden)
		if err := json.NewEncoder(w).Encode(PATCHResp{
			Error: "You do not have access to this board.",
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	}

	// Update the task and subtasks in the database.
	if err = h.taskUpdater.Update(id, taskTable.NewUpRecord(
		reqBody.Title, reqBody.Description, subtaskRecords,
	)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}
}
