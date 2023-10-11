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

// DELETEHandler is an api.MethodHandler that can be used to handle DELETE
// requests sent to the task route.
type DELETEHandler struct {
	idValidator       api.StringValidator
	taskSelector      dbaccess.Selector[taskTable.Record]
	columnSelector    dbaccess.Selector[columnTable.Record]
	userBoardSelector dbaccess.RelSelector[bool]
	taskDeleter       dbaccess.Deleter
	log               pkgLog.Errorer
}

// NewDELETEHandler creates and returns a new DELETEHandler.
func NewDELETEHandler(
	idValidator api.StringValidator,
	taskSelector dbaccess.Selector[taskTable.Record],
	columnSelector dbaccess.Selector[columnTable.Record],
	userBoardSelector dbaccess.RelSelector[bool],
	taskDeleter dbaccess.Deleter,
	log pkgLog.Errorer,
) DELETEHandler {
	return DELETEHandler{
		idValidator:       idValidator,
		taskSelector:      taskSelector,
		columnSelector:    columnSelector,
		userBoardSelector: userBoardSelector,
		taskDeleter:       taskDeleter,
		log:               log,
	}
}

// Handle handles the DELETE requests sent to the task route.
func (h DELETEHandler) Handle(
	w http.ResponseWriter, r *http.Request, username string,
) {
	// Read and validate task ID.
	id := r.URL.Query().Get("id")
	if err := h.idValidator.Validate(id); errors.Is(err, api.ErrStrEmpty) {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(ResBody{
			Error: "Task ID cannot be empty.",
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
			return
		}
	} else if errors.Is(err, api.ErrStrNotInt) {
		w.WriteHeader(http.StatusBadRequest)
		if err = json.NewEncoder(w).Encode(ResBody{
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
		if err = json.NewEncoder(w).Encode(ResBody{
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

	// Authorize the user.
	isAdmin, err := h.userBoardSelector.Select(username, strconv.Itoa(column.BoardID))
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusForbidden)
		if err = json.NewEncoder(w).Encode(ResBody{
			Error: "You do not have access to this board.",
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
	if !isAdmin {
		w.WriteHeader(http.StatusForbidden)
		if err = json.NewEncoder(w).Encode(ResBody{
			Error: "Only board admins can delete tasks.",
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
			return
		}
	}

	// Delete the record from task table that has the given ID.
	if err = h.taskDeleter.Delete(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}
}
