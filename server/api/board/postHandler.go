package board

import (
	"encoding/json"
	"net/http"

	"server/db"
	"server/relay"
)

// POSTHandler can be used to handle the POST requests sent to the board
// endpoint.
type POSTHandler struct {
	userBoardCounter db.Counter
	boardInserter    db.Inserter[db.Board]
}

// NewPOSTHandler creates and returns a new POSTHandler.
func NewPOSTHandler(
	userBoardCounter db.Counter,
	boardInserter db.Inserter[db.Board],
) POSTHandler {
	return POSTHandler{
		userBoardCounter: userBoardCounter,
		boardInserter:    boardInserter,
	}
}

// Handle handles the POST requests sent to the board endpoint.
func (h POSTHandler) Handle(
	w http.ResponseWriter, r *http.Request, sub string,
) {
	// Read the request body.
	reqBody := ReqBody{}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		relay.ServerErr(w, err.Error())
		return
	}

	// Validate board name.
	if reqBody.Name == "" {
		relay.ClientErr(
			w, http.StatusBadRequest, ResBody{Error: errNameEmpty},
		)
		return
	}
	if len(reqBody.Name) >= maxNameLength {
		relay.ClientErr(
			w, http.StatusBadRequest, ResBody{Error: errNameTooLong},
		)
		return
	}

	// Validate that the user has less than 3 boards. This is done to limit the
	// resources used by this demo app.
	if boardCount := h.userBoardCounter.Count(sub); boardCount >= maxBoards {
		relay.ClientErr(w, http.StatusBadRequest, ResBody{Error: errMaxBoards})
		return
	}

	// Create a new board.
	if err := h.boardInserter.Insert(
		db.NewBoard(reqBody.Name, sub),
	); err != nil {
		relay.ServerErr(w, err.Error())
		return
	}
}

const (
	// maxBoards is the amount of boards that each user is allowed to own (i.e.
	// be the admin of).
	maxBoards = 3

	// maxNameLength is the maximum amount of characters that a board name can
	// have.
	maxNameLength = 35

	// errMaxBoards is the error message returned from the handler when the user
	// already owns the maximum amount of boards allowed per user.
	errMaxBoards = "You have already created the maximum amount of boards all" +
		"owed per user. Please delete one of your boards to create a new one."

	// errNameEmpty is the error message returned from the handler when the
	// received board name is empty.
	errNameEmpty = "Board name cannot be empty."

	// errNameTooLong is the error message returned from the handler when the
	// received board name is too long.
	errNameTooLong = "Board name cannot be longer than 35 characters."
)
