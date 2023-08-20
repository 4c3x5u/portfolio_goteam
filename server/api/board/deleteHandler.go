package board

import (
	"database/sql"
	"errors"
	"net/http"

	"server/dbaccess"
	pkgLog "server/log"
)

// DELETEHandler is an api.MethodHandler that can be used to handle DELETE board
// requests.
type DELETEHandler struct {
	validator         DELETEReqValidator
	userBoardSelector dbaccess.RelSelector[bool]
	boardDeleter      dbaccess.Deleter
	log               pkgLog.Errorer
}

// NewDELETEHandler creates and returns a new DELETEHandler.
func NewDELETEHandler(
	validator DELETEReqValidator,
	userBoardSelector dbaccess.RelSelector[bool],
	boardDeleter dbaccess.Deleter,
	log pkgLog.Errorer,
) DELETEHandler {
	return DELETEHandler{
		validator:         validator,
		userBoardSelector: userBoardSelector,
		boardDeleter:      boardDeleter,
		log:               log,
	}
}

// Handle handles the DELETE requests sent to the board route.
func (h DELETEHandler) Handle(
	w http.ResponseWriter, r *http.Request, username string,
) {
	// Get id query parameter. That's our board ID.
	boardID := r.URL.Query().Get("id")

	if ok := h.validator.Validate(boardID); !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Validate that the user making the request is the admin of the board to be
	// deleted.
	if isAdmin, err := h.userBoardSelector.Select(
		username, boardID,
	); err != nil && err != sql.ErrNoRows {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if !isAdmin {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Delete the board.
	if err := h.boardDeleter.Delete(boardID); err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// All went well. Return 200.
	w.WriteHeader(http.StatusOK)
}
