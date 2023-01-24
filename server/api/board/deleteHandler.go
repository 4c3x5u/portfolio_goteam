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
}

// NewDELETEHandler creates and returns a new DELETEHandler.
func NewDELETEHandler(userBoardSelector db.RelSelector[bool]) DELETEHandler {
	return DELETEHandler{userBoardSelector: userBoardSelector}
}

// Handle handles the DELETE requests sent to the board endpoint.
func (h DELETEHandler) Handle(
	w http.ResponseWriter, r *http.Request, sub string,
) {
	// Get id query parameter. That's our board ID.
	boardID := r.URL.Query().Get("id")

	// Secelt isAdmin from user-board relationship table for column
	// that matches the sender's username (i.e. sub) AND the boardID
	// passed through the query parameter.
	isAdmin, err := h.userBoardSelector.Select(sub, boardID)
	if err == sql.ErrNoRows {
		// no rows found (404)
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		// other error (500)
		relay.ServerErr(w, err.Error())
		return
	}
	if isAdmin == false {
		// user not admin (401)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// todo: make 200 when deletehandler work is done
	w.WriteHeader(http.StatusNotImplemented)
	return
}
