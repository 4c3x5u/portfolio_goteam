package board

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"server/db"
	"server/log"
)

// POSTHandler is an api.MethodHandler that can be used to handle POST board
// requests.
type POSTHandler struct {
	validator        POSTReqValidator
	userBoardCounter db.Counter
	boardInserter    db.Inserter[db.Board]
	logger           log.Logger
}

// NewPOSTHandler creates and returns a new POSTHandler.
func NewPOSTHandler(
	validator POSTReqValidator,
	userBoardCounter db.Counter,
	boardInserter db.Inserter[db.Board],
	logger log.Logger,
) POSTHandler {
	return POSTHandler{
		validator:        validator,
		userBoardCounter: userBoardCounter,
		boardInserter:    boardInserter,
		logger:           logger,
	}
}

// Handle handles the POST requests sent to the board route.
func (h POSTHandler) Handle(
	w http.ResponseWriter, r *http.Request, username string,
) {
	// Read and validate request body.
	reqBody := POSTReqBody{}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		h.logger.Log(log.LevelError, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if errMsg := h.validator.Validate(reqBody); errMsg != "" {
		w.WriteHeader(http.StatusBadRequest)
		if encodeErr := json.NewEncoder(w).Encode(
			POSTResBody{Error: errMsg},
		); encodeErr != nil {
			h.logger.Log(log.LevelError, encodeErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// Validate that the user has less than 3 boards. This is done to limit the
	// resources used by this demo app.
	if boardCount, err := h.userBoardCounter.Count(
		username,
	); err != nil && err != sql.ErrNoRows {
		// sql.ErrNoRows is OK here. It just means the user hasn't created any
		// boards yet.
		h.logger.Log(log.LevelError, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if boardCount >= 3 {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(
			POSTResBody{
				Error: "You have already created the maximum amount of " +
					"boards allowed per user. Please delete one of your " +
					"boards to create a new one.",
			},
		); err != nil {
			h.logger.Log(log.LevelError, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
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
}
