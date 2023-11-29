package subtask

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
	subtaskTable "github.com/kxplxn/goteam/server/dbaccess/subtask"
	taskTable "github.com/kxplxn/goteam/server/dbaccess/task"
	userTable "github.com/kxplxn/goteam/server/dbaccess/user"
	pkgLog "github.com/kxplxn/goteam/server/log"
)

// PATCHReq defines the body of PATCH subtask requests.
type PATCHReq struct {
	IsDone bool `json:"done"`
}

// PATCHResp defines the body of PATCH subtask response.
type PATCHResp struct {
	Error string `json:"error"`
}

// PATCHHandler is an api.MethodHandler that can be used to handle PATCH
// requests sent to the subtask route.
type PATCHHandler struct {
	userSelector    dbaccess.Selector[userTable.Record]
	idValidator     api.StringValidator
	subtaskSelector dbaccess.Selector[subtaskTable.Record]
	taskSelector    dbaccess.Selector[taskTable.Record]
	columnSelector  dbaccess.Selector[columnTable.Record]
	boardSelector   dbaccess.Selector[boardTable.Record]
	subtaskUpdater  dbaccess.Updater[subtaskTable.UpRecord]
	log             pkgLog.Errorer
}

// NewPATCHHandler creates and returns a new PATCHandler.
func NewPATCHHandler(
	userSelector dbaccess.Selector[userTable.Record],
	idValidator api.StringValidator,
	subtaskSelector dbaccess.Selector[subtaskTable.Record],
	taskSelector dbaccess.Selector[taskTable.Record],
	columnSelector dbaccess.Selector[columnTable.Record],
	boardSelector dbaccess.Selector[boardTable.Record],
	subtaskUpdater dbaccess.Updater[subtaskTable.UpRecord],
	log pkgLog.Errorer,
) PATCHHandler {
	return PATCHHandler{
		userSelector:    userSelector,
		idValidator:     idValidator,
		subtaskSelector: subtaskSelector,
		taskSelector:    taskSelector,
		columnSelector:  columnSelector,
		boardSelector:   boardSelector,
		subtaskUpdater:  subtaskUpdater,
		log:             log,
	}
}

// Handle handles the PATCH requests sent to the subtask route.
func (h PATCHHandler) Handle(
	w http.ResponseWriter, r *http.Request, username string,
) {
	// Validate that the user is a team admin.
	user, err := h.userSelector.Select(username)
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusUnauthorized)
		if err := json.NewEncoder(w).Encode(
			PATCHResp{Error: "Username is not recognised."},
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
	if !user.IsAdmin {
		w.WriteHeader(http.StatusForbidden)
		if err := json.NewEncoder(w).Encode(
			PATCHResp{Error: "Only team admins can edit subtasks."},
		); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	}

	// Read and validate subtask ID.
	id := r.URL.Query().Get("id")
	if err := h.idValidator.Validate(id); errors.Is(err, api.ErrStrEmpty) {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(
			PATCHResp{Error: "Subtask ID cannot be empty."},
		); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	} else if errors.Is(err, api.ErrStrNotInt) {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(
			PATCHResp{Error: "Subtask ID must be an integer."},
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
			PATCHResp{Error: "Subtask not found."},
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

	// Retrieve board to compare its TeamID to the user's TeamID.
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

	// Update subtask based on request body.
	var reqBody PATCHReq
	if err = json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}
	if err = h.subtaskUpdater.Update(
		id, subtaskTable.NewUpRecord(reqBody.IsDone),
	); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}
}
