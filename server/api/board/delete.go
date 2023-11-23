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
	validator     api.StringValidator
	userSelector  dbaccess.Selector[userTable.Record]
	boardSelector dbaccess.Selector[boardTable.Record]
	boardDeleter  dbaccess.Deleter
	log           pkgLog.Errorer
}

// NewDELETEHandler creates and returns a new DELETEHandler.
func NewDELETEHandler(
	validator api.StringValidator,
	userSelector dbaccess.Selector[userTable.Record],
	boardSelector dbaccess.Selector[boardTable.Record],
	boardDeleter dbaccess.Deleter,
	log pkgLog.Errorer,
) DELETEHandler {
	return DELETEHandler{
		validator:     validator,
		userSelector:  userSelector,
		boardSelector: boardSelector,
		boardDeleter:  boardDeleter,
		log:           log,
	}
}

// Handle handles the DELETE requests sent to the board route.
func (h DELETEHandler) Handle(
	w http.ResponseWriter, r *http.Request, username string,
) {
	// Get id query parameter. That's our board ID.
	boardID := r.URL.Query().Get("id")

	// Validate board ID.
	if err := h.validator.Validate(boardID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Validate that the user making the request is the admin of the board to be
	// deleted.
	user, err := h.userSelector.Select(
		username,
	)
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
