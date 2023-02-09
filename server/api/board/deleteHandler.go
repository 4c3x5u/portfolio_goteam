package board

import (
	"database/sql"
	"net/http"

	"server/db"
	"server/log"
)

// DELETEHandler is a MethodHandler that is intended to handle DELETE requests
// sent to the board route.
type DELETEHandler struct {
	validator         DELETEReqValidator
	userBoardSelector db.RelSelector[bool]
	boardDeleter      db.Deleter
	logger            log.Logger
}

// NewDELETEHandler creates and returns a new DELETEHandler.
func NewDELETEHandler(
	validator DELETEReqValidator,
	userBoardSelector db.RelSelector[bool],
	boardDeleter db.Deleter,
	logger log.Logger,
) DELETEHandler {
	return DELETEHandler{
		validator:         validator,
		userBoardSelector: userBoardSelector,
		boardDeleter:      boardDeleter,
		logger:            logger,
	}
}

// Handle handles the DELETE requests sent to the board route.
func (h DELETEHandler) Handle(
	w http.ResponseWriter, r *http.Request, username string,
) {
	// Get id query parameter. That's our board ID.
	boardID, ok := h.validator.Validate(r.URL.Query())
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Validate that the user making the request is the admin of the board to be
	// deleted.
	isAdmin, err := h.userBoardSelector.Select(username, boardID)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		h.logger.Log(log.LevelError, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !isAdmin {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Delete the board.
	if err = h.boardDeleter.Delete(boardID); err != nil {
		h.logger.Log(log.LevelError, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// All went well. Return 200.
	w.WriteHeader(http.StatusOK)
}
