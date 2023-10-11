package column

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/kxplxn/goteam/server/api"
	"github.com/kxplxn/goteam/server/dbaccess"
	columnTable "github.com/kxplxn/goteam/server/dbaccess/column"
	pkgLog "github.com/kxplxn/goteam/server/log"
)

// PATCHHandler is an api.MethodHandler that can be used to handle PATCH
// requests sent to the column route.
type PATCHHandler struct {
	idValidator       api.StringValidator
	columnSelector    dbaccess.Selector[columnTable.Record]
	userBoardSelector dbaccess.RelSelector[bool]
	columnUpdater     dbaccess.Updater[[]columnTable.Task]
	log               pkgLog.Errorer
}

// NewPATCHHandler creates and returns a new PATCHHandler.
func NewPATCHHandler(
	idValidator api.StringValidator,
	columnSelector dbaccess.Selector[columnTable.Record],
	userBoardSelector dbaccess.RelSelector[bool],
	columnUpdater dbaccess.Updater[[]columnTable.Task],
	log pkgLog.Errorer,
) PATCHHandler {
	return PATCHHandler{
		idValidator:       idValidator,
		columnSelector:    columnSelector,
		userBoardSelector: userBoardSelector,
		columnUpdater:     columnUpdater,
		log:               log,
	}
}

// Handle handles the PATCH requests sent to the column route.
func (h PATCHHandler) Handle(
	w http.ResponseWriter, r *http.Request, username string,
) {
	// Get and validate the column ID.
	columnID := r.URL.Query().Get("id")
	if err := h.idValidator.Validate(columnID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err = json.NewEncoder(w).Encode(
			ResBody{Error: err.Error()},
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
		w.WriteHeader(http.StatusBadRequest)
		if err = json.NewEncoder(w).Encode(
			ResBody{Error: "Column not found."},
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

	// Check whether the user has the right to edit this column.
	if isAdmin, err := h.userBoardSelector.Select(
		username, strconv.Itoa(column.BoardID),
	); errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusUnauthorized)
		if err = json.NewEncoder(w).Encode(
			ResBody{Error: "You do not have access to this board."},
		); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	} else if !isAdmin {
		w.WriteHeader(http.StatusUnauthorized)
		if err = json.NewEncoder(w).Encode(
			ResBody{Error: "Only board admins can move tasks."},
		); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	}

	// Decode request body and map it into tasks.
	var reqBody ReqBody
	if err = json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}
	var tasks []columnTable.Task
	for _, t := range reqBody {
		tasks = append(tasks, columnTable.Task{ID: t.ID, Order: t.Order})
	}

	// Update task records in the database using column ID and order from tasks.
	if err = h.columnUpdater.Update(
		columnID, tasks,
	); errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		if err = json.NewEncoder(w).Encode(
			ResBody{Error: "Task not found."},
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
