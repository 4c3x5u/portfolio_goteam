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
	log               pkgLog.Errorer
}

// NewDELETEHandler creates and returns a new DELETEHandler.
func NewDELETEHandler(
	idValidator api.StringValidator,
	taskSelector dbaccess.Selector[taskTable.Record],
	columnSelector dbaccess.Selector[columnTable.Record],
	userBoardSelector dbaccess.RelSelector[bool],
	log pkgLog.Errorer,
) DELETEHandler {
	return DELETEHandler{
		idValidator:       idValidator,
		taskSelector:      taskSelector,
		columnSelector:    columnSelector,
		userBoardSelector: userBoardSelector,
		log:               log,
	}
}

// Handle handles the DELETE requests sent to the task route.
func (h DELETEHandler) Handle(
	w http.ResponseWriter, r *http.Request, username string,
) {
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

	column, err := h.columnSelector.Select(strconv.Itoa(task.ColumnID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}

	_, err = h.userBoardSelector.Select(username, strconv.Itoa(column.BoardID))
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
}
