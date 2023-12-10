package task

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
	taskTable "github.com/kxplxn/goteam/pkg/dbaccess/task"
	pkgLog "github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/pkg/token"
)

// DeleteResp defines the body of DELETE task responses.
type DeleteResp struct {
	Error string `json:"error"`
}

// DeleteHandler is an api.MethodHandler that can be used to handle DELETE
// requests made to the task route.
type DeleteHandler struct {
	decodeAuth     token.DecodeFunc[token.Auth]
	idValidator    api.StringValidator
	taskSelector   dbaccess.Selector[taskTable.Record]
	columnSelector dbaccess.Selector[columnTable.Record]
	boardSelector  dbaccess.Selector[boardTable.Record]
	taskDeleter    dbaccess.Deleter
	log            pkgLog.Errorer
}

// NewDeleteHandler creates and returns a new DELETEHandler.
func NewDeleteHandler(
	decodeAuth token.DecodeFunc[token.Auth],
	idValidator api.StringValidator,
	taskSelector dbaccess.Selector[taskTable.Record],
	columnSelector dbaccess.Selector[columnTable.Record],
	boardSelector dbaccess.Selector[boardTable.Record],
	taskDeleter dbaccess.Deleter,
	log pkgLog.Errorer,
) DeleteHandler {
	return DeleteHandler{
		decodeAuth:     decodeAuth,
		idValidator:    idValidator,
		taskSelector:   taskSelector,
		columnSelector: columnSelector,
		boardSelector:  boardSelector,
		taskDeleter:    taskDeleter,
		log:            log,
	}
}

// Handle handles the DELETE requests sent to the task route.
func (h DeleteHandler) Handle(
	w http.ResponseWriter, r *http.Request, username string,
) {
	// get auth token
	ckAuth, err := r.Cookie(token.AuthName)
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusUnauthorized)
		if encodeErr := json.NewEncoder(w).Encode(PostResp{
			Error: "Auth token not found.",
		}); encodeErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}

	// decode auth token and validate user is admin
	user, err := h.decodeAuth(ckAuth.Value)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		if err = json.NewEncoder(w).Encode(DeleteResp{
			Error: "Invalid auth token.",
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
			return
		}
	}
	if !user.IsAdmin {
		w.WriteHeader(http.StatusForbidden)
		if err = json.NewEncoder(w).Encode(DeleteResp{
			Error: "Only board admins can delete tasks.",
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
			return
		}
	}

	// Read and validate task ID.
	id := r.URL.Query().Get("id")
	if err := h.idValidator.Validate(id); errors.Is(err, api.ErrEmpty) {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(DeleteResp{
			Error: "Task ID cannot be empty.",
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
			return
		}
	} else if errors.Is(err, api.ErrNotInt) {
		w.WriteHeader(http.StatusBadRequest)
		if err = json.NewEncoder(w).Encode(DeleteResp{
			Error: "Task ID must be an integer.",
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
			return
		}
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}

	// Select task from the database to access its column ID.
	task, err := h.taskSelector.Select(id)
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		if err = json.NewEncoder(w).Encode(DeleteResp{
			Error: "Task not found.",
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
			return
		}
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}

	// Select column from the database to access its board ID.
	column, err := h.columnSelector.Select(strconv.Itoa(task.ColumnID))
	if err != nil {
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
	if user.TeamID != strconv.Itoa(board.TeamID) {
		w.WriteHeader(http.StatusForbidden)
		if encodeErr := json.NewEncoder(w).Encode(DeleteResp{
			Error: "You do not have access to this board.",
		}); encodeErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	}

	// Delete the record from task table that has the given ID.
	if err = h.taskDeleter.Delete(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}
}
