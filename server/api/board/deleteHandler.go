package board

import (
	"database/sql"
	"net/http"

	"server/db"
	"server/relay"
)

// DELETEHandler can be used to handle the DELETE requests sent to the board
// endpoint.
type DELETEHandler struct {
	userBoardSelector db.RelSelector[bool]
	boardDeleter      db.Deleter
}

// NewDELETEHandler creates and returns a new DELETEHandler.
func NewDELETEHandler(
	userBoardSelector db.RelSelector[bool],
	boardDeleter db.Deleter,
) DELETEHandler {
	return DELETEHandler{
		userBoardSelector: userBoardSelector,
		boardDeleter:      boardDeleter,
	}
}

// Handle handles the DELETE requests sent to the board endpoint.
func (h DELETEHandler) Handle(
	w http.ResponseWriter, r *http.Request, username string,
) {
	// Get id query parameter. That's our board ID.
	boardID := r.URL.Query().Get("id")

	// Validate that the user making the request is the admin of the board to be
	// deleted.
	isAdmin, err := h.userBoardSelector.Select(username, boardID)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		relay.ServerErr(w, err.Error())
		return
	}
	if isAdmin == false {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Delete the board.
	if err = h.boardDeleter.Delete(boardID); err != nil {
		relay.ServerErr(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
