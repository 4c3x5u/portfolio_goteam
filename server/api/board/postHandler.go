package board

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"server/db"
	"server/log"
)

// POSTHandler can be used to handle the POST requests sent to the board
// endpoint.
type POSTHandler struct {
	userBoardCounter db.Counter
	boardInserter    db.Inserter[db.Board]
	logger           log.Logger
}

// NewPOSTHandler creates and returns a new POSTHandler.
func NewPostHandler(
	userBoardCounter db.Counter,
	boardInserter db.Inserter[db.Board],
	logger log.Logger,
) POSTHandler {
	return POSTHandler{
		userBoardCounter: userBoardCounter,
		boardInserter:    boardInserter,
		logger:           logger,
	}
}

// Handle handles the POST requests sent to the board endpoint.
func (h POSTHandler) Handle(
	w http.ResponseWriter, r *http.Request, username string,
) {
	// Read the request body.
	reqBody := ReqBody{}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		h.logger.Log(log.LevelError, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Validate board name.
	if reqBody.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(
			ResBody{Error: errNameEmpty},
		); err != nil {
			h.logger.Log(log.LevelError, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	if len(reqBody.Name) >= maxNameLength {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(
			ResBody{Error: errNameTooLong},
		); err != nil {
			h.logger.Log(log.LevelError, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	// Validate that the user has less than 3 boards. This is done to limit the
	// resources used by this demo app.
	if boardCount, err := h.userBoardCounter.Count(
		username,
	); err != nil && err != sql.ErrNoRows {
		h.logger.Log(log.LevelError, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if boardCount >= maxBoards {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(
			ResBody{Error: errMaxBoards},
		); err != nil {
			h.logger.Log(log.LevelError, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	// Create a new board.
	if err := h.boardInserter.Insert(
		db.NewBoard(reqBody.Name, username),
	); err != nil {
		h.logger.Log(log.LevelError, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// All went well. Return 200.
	w.WriteHeader(http.StatusOK)
	return
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
