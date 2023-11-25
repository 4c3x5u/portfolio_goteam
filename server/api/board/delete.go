package board

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/kxplxn/goteam/server/api"
	"github.com/kxplxn/goteam/server/dbaccess"
	boardTable "github.com/kxplxn/goteam/server/dbaccess/board"
	userTable "github.com/kxplxn/goteam/server/dbaccess/user"
	pkgLog "github.com/kxplxn/goteam/server/log"
)

// DELETEHandler is an api.MethodHandler that can be used to handle DELETE board
// requests.
type DELETEHandler struct {
	userSelector  dbaccess.Selector[userTable.Record]
	validator     api.StringValidator
	boardSelector dbaccess.Selector[boardTable.Record]
	boardDeleter  dbaccess.Deleter
	log           pkgLog.Errorer
}

// NewDELETEHandler creates and returns a new DELETEHandler.
func NewDELETEHandler(
	userSelector dbaccess.Selector[userTable.Record],
	validator api.StringValidator,
	boardSelector dbaccess.Selector[boardTable.Record],
	boardDeleter dbaccess.Deleter,
	log pkgLog.Errorer,
) DELETEHandler {
	return DELETEHandler{
		userSelector:  userSelector,
		validator:     validator,
		boardSelector: boardSelector,
		boardDeleter:  boardDeleter,
		log:           log,
	}
}

// Handle handles the DELETE requests sent to the board route.
func (h DELETEHandler) Handle(
	w http.ResponseWriter, r *http.Request, username string,
) {
	// Validate that the user is a team admin.
	user, err := h.userSelector.Select(username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if !user.IsAdmin {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Get id query parameter. That's our board ID.
	boardID := r.URL.Query().Get("id")

	// Validate board ID.
	if err := h.validator.Validate(boardID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Select board to validate that it belongs to the team that the user is the
	// admin of.
	board, err := h.boardSelector.Select(boardID)
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if board.TeamID != user.TeamID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Delete the board.
	if err = h.boardDeleter.Delete(boardID); err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// All went well. Return 200.
	w.WriteHeader(http.StatusOK)
}
