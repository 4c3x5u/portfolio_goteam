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

// POSTReqBody defines the request body for requests handled by POSTHandler.
type POSTReqBody struct {
	ColumnID      int      `json:"column"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	SubtaskTitles []string `json:"subtasks"`
}

// POSTHandler is an api.MethodHandler that can be used to handle POST requests
// sent to the task route.
type POSTHandler struct {
	userSelector          dbaccess.Selector[userTable.Record]
	taskTitleValidator    api.StringValidator
	subtaskTitleValidator api.StringValidator
	columnSelector        dbaccess.Selector[columnTable.Record]
	boardSelector         dbaccess.Selector[boardTable.Record]
	taskInserter          dbaccess.Inserter[taskTable.InRecord]
	log                   pkgLog.Errorer
}

// NewPOSTHandler creates and returns a new POSTHandler.
func NewPOSTHandler(
	userSelector dbaccess.Selector[userTable.Record],
	taskTitleValidator api.StringValidator,
	subtaskTitleValidator api.StringValidator,
	columnSelector dbaccess.Selector[columnTable.Record],
	boardSelector dbaccess.Selector[boardTable.Record],
	taskInserter dbaccess.Inserter[taskTable.InRecord],
	log pkgLog.Errorer,
) *POSTHandler {
	return &POSTHandler{
		userSelector:          userSelector,
		taskTitleValidator:    taskTitleValidator,
		subtaskTitleValidator: subtaskTitleValidator,
		columnSelector:        columnSelector,
		boardSelector:         boardSelector,
		taskInserter:          taskInserter,
		log:                   log,
	}
}

// Handle handles the POST requests sent to the task route.
func (h *POSTHandler) Handle(
	w http.ResponseWriter, r *http.Request, username string,
) {
	// Check if the user is a team admin.
	user, err := h.userSelector.Select(username)
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusUnauthorized)
		if encodeErr := json.NewEncoder(w).Encode(ResBody{
			Error: "Username is not recognised.",
		}); encodeErr != nil {
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
		if encodeErr := json.NewEncoder(w).Encode(ResBody{
			Error: "Only team admins can create tasks.",
		}); encodeErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	}

	var reqBody POSTReqBody
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
		if err = json.NewEncoder(w).Encode(ResBody{
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
			if err = json.NewEncoder(w).Encode(ResBody{
				Error: errMsg,
			}); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				h.log.Error(err.Error())
			}
			return
		}
	}

	// Get the column from the database with the task's column ID.
	column, err := h.columnSelector.Select(
		strconv.Itoa(reqBody.ColumnID),
	)
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		if encodeErr := json.NewEncoder(w).Encode(ResBody{
			Error: "Column not found.",
		}); encodeErr != nil {
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

	// Check if the board belongs to the team that the user is the admin of.
	board, err := h.boardSelector.Select(strconv.Itoa(column.BoardID))
	if err != nil {
		// Since boardID is used from a column retrieved from the database,
		// any error selecting board including ErrNoRows is a 500.
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}
	if board.TeamID != user.TeamID {
		w.WriteHeader(http.StatusForbidden)
		if encodeErr := json.NewEncoder(w).Encode(ResBody{
			Error: "You do not have access to this board.",
		}); encodeErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	}

	// Insert task and subtasks into the database.
	if err = h.taskInserter.Insert(taskTable.NewInRecord(
		reqBody.ColumnID,
		reqBody.Title,
		reqBody.Description,
		reqBody.SubtaskTitles,
	)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}
}
