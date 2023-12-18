package board

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/kxplxn/goteam/internal/api"
	"github.com/kxplxn/goteam/pkg/legacydb"
	boardTable "github.com/kxplxn/goteam/pkg/legacydb/board"
	pkgLog "github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/pkg/token"
)

// DELETEHandler is an api.MethodHandler that can be used to handle DELETE board
// requests.
type DELETEHandler struct {
	decodeAuth    token.DecodeFunc[token.Auth]
	validator     api.StringValidator
	boardSelector legacydb.Selector[boardTable.Record]
	boardDeleter  legacydb.Deleter
	log           pkgLog.Errorer
}

// NewDELETEHandler creates and returns a new DELETEHandler.
func NewDELETEHandler(
	decodeAuth token.DecodeFunc[token.Auth],
	validator api.StringValidator,
	boardSelector legacydb.Selector[boardTable.Record],
	boardDeleter legacydb.Deleter,
	log pkgLog.Errorer,
) DELETEHandler {
	return DELETEHandler{
		decodeAuth:    decodeAuth,
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
	ckAuth, err := r.Cookie(token.AuthName)
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}

	// decode auth token
	auth, err := h.decodeAuth(ckAuth.Value)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// validate user is admin
	if !auth.IsAdmin {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// get and validate board ID
	boardID := r.URL.Query().Get("id")
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
	if strconv.Itoa(board.TeamID) != auth.TeamID {
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
