package tasks

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/kxplxn/goteam/internal/api"
	"github.com/kxplxn/goteam/pkg/dbaccess"
	boardTable "github.com/kxplxn/goteam/pkg/dbaccess/board"
	columnTable "github.com/kxplxn/goteam/pkg/dbaccess/column"
	userTable "github.com/kxplxn/goteam/pkg/dbaccess/user"
	pkgLog "github.com/kxplxn/goteam/pkg/log"
)

// PatchReq defines body of PATCH tasks requests.
type PatchReq []Task

// Task represents an element in PatchReq.
type Task struct {
	ID    int `json:"id"`
	Order int `json:"order"`
}

// PatchResp defines the body for PATCH column responses.
type PatchResp struct {
	Error string `json:"error"`
}

// PatchHandler is an api.MethodHandler that can be used to handle PATCH
// requests sent to the tasks route.
type PatchHandler struct {
	userSelector   dbaccess.Selector[userTable.Record]
	idValidator    api.StringValidator
	columnSelector dbaccess.Selector[columnTable.Record]
	boardSelector  dbaccess.Selector[boardTable.Record]
	columnUpdater  dbaccess.Updater[[]columnTable.Task]
	log            pkgLog.Errorer
}

// NewPATCHHandler creates and returns a new PATCHHandler.
func NewPATCHHandler(
	userSelector dbaccess.Selector[userTable.Record],
	idValidator api.StringValidator,
	columnSelector dbaccess.Selector[columnTable.Record],
	boardSelector dbaccess.Selector[boardTable.Record],
	columnUpdater dbaccess.Updater[[]columnTable.Task],
	log pkgLog.Errorer,
) PatchHandler {
	return PatchHandler{
		userSelector:   userSelector,
		idValidator:    idValidator,
		columnSelector: columnSelector,
		boardSelector:  boardSelector,
		columnUpdater:  columnUpdater,
		log:            log,
	}
}

// Handle handles the PATCH requests sent to the tasks route.
func (h PatchHandler) Handle(
	w http.ResponseWriter, r *http.Request, username string,
) {
	// Validate that the user is a team admin.
	user, err := h.userSelector.Select(username)
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusUnauthorized)
		if err = json.NewEncoder(w).Encode(
			PatchResp{Error: "Username is not recognised."},
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
		if err = json.NewEncoder(w).Encode(
			PatchResp{Error: "Only team admins can move tasks."},
		); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	}

	// Get and validate the column ID.
	columnID := r.URL.Query().Get("id")
	if err := h.idValidator.Validate(columnID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err = json.NewEncoder(w).Encode(
			PatchResp{Error: err.Error()},
		); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	}

	// Retrieve the column from the database so that we find out its board ID to
	// validate that the user has the right to edit it.
	column, err := h.columnSelector.Select(columnID)
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		if err = json.NewEncoder(w).Encode(
			PatchResp{Error: "Column not found."},
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

	// Validate that the column's board belongs to the team that the user is the
	// admin of.
	board, err := h.boardSelector.Select(strconv.Itoa(column.BoardID))
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		if err = json.NewEncoder(w).Encode(
			PatchResp{Error: "Board not found."},
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
	if board.TeamID != user.TeamID {
		w.WriteHeader(http.StatusForbidden)
		if err = json.NewEncoder(w).Encode(
			PatchResp{Error: "You do not have access to this board."},
		); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	}

	// Decode request body and map it into tasks.
	var req PatchReq
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}
	var tasks []columnTable.Task
	for _, t := range req {
		tasks = append(tasks, columnTable.Task{ID: t.ID, Order: t.Order})
	}

	// Update task records in the database using column ID and order from tasks.
	if err = h.columnUpdater.Update(
		columnID, tasks,
	); errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		if err = json.NewEncoder(w).Encode(
			PatchResp{Error: "Task not found."},
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

	// All went well. Return 200.
	w.WriteHeader(http.StatusOK)
}
