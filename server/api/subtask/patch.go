package subtask

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/kxplxn/goteam/server/api"
	"github.com/kxplxn/goteam/server/dbaccess"
	columnTable "github.com/kxplxn/goteam/server/dbaccess/column"
	subtaskTable "github.com/kxplxn/goteam/server/dbaccess/subtask"
	taskTable "github.com/kxplxn/goteam/server/dbaccess/task"
	pkgLog "github.com/kxplxn/goteam/server/log"
)

// ResBody defines the response body for requests handled by PATCHHandler.
type ResBody struct {
	Error string `json:"error"`
}

// PATCHHandler is an api.MethodHandler that can be used to handle PATCH
// requests sent to the subtask route.
type PATCHHandler struct {
	idValidator       api.StringValidator
	subtaskSelector   dbaccess.Selector[subtaskTable.Record]
	taskSelector      dbaccess.Selector[taskTable.Record]
	columnSelector    dbaccess.Selector[columnTable.Record]
	userBoardSelector dbaccess.RelSelector[bool]
	log               pkgLog.Errorer
}

// NewPATCHHandler creates and returns a new PATCHandler.
func NewPATCHHandler(
	idValidator api.StringValidator,
	subtaskSelector dbaccess.Selector[subtaskTable.Record],
	taskSelector dbaccess.Selector[taskTable.Record],
	columnSelector dbaccess.Selector[columnTable.Record],
	userBoardSelector dbaccess.RelSelector[bool],
	log pkgLog.Errorer,
) PATCHHandler {
	return PATCHHandler{
		idValidator:       idValidator,
		subtaskSelector:   subtaskSelector,
		taskSelector:      taskSelector,
		columnSelector:    columnSelector,
		userBoardSelector: userBoardSelector,
		log:               log,
	}
}

// Handle handles the PATCH requests sent to the subtask route.
func (h PATCHHandler) Handle(
	w http.ResponseWriter, r *http.Request, username string,
) {
	// Read and validate subtask ID.
	id := r.URL.Query().Get("id")
	if err := h.idValidator.Validate(id); errors.Is(err, api.ErrStrEmpty) {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(
			ResBody{Error: "Subtask ID cannot be empty."},
		); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	} else if errors.Is(err, api.ErrStrNotInt) {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(
			ResBody{Error: "Subtask ID must be an integer."},
		); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}

	// Retrieve subtask to access task ID.
	subtask, err := h.subtaskSelector.Select(id)
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(
			ResBody{Error: "Subtask not found."},
		); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}

	// Retrieve task to access column ID.
	task, err := h.taskSelector.Select(strconv.Itoa(subtask.TaskID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}

	// Retrieve column to access board ID.
	column, err := h.columnSelector.Select(strconv.Itoa(task.ColumnID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}

	// Authorize the user.
	isAdmin, err := h.userBoardSelector.Select(
		username, strconv.Itoa(column.BoardID),
	)
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusForbidden)
		if err := json.NewEncoder(w).Encode(
			ResBody{Error: "You do not have access to this board."},
		); err != nil {
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
		if err := json.NewEncoder(w).Encode(
			ResBody{Error: "Only board admins can edit subtasks."},
		); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	}
}
