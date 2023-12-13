package board

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/kxplxn/goteam/internal/api"
	"github.com/kxplxn/goteam/pkg/legacydb"
	boardTable "github.com/kxplxn/goteam/pkg/legacydb/board"
	teamTable "github.com/kxplxn/goteam/pkg/legacydb/team"
	userTable "github.com/kxplxn/goteam/pkg/legacydb/user"
	pkgLog "github.com/kxplxn/goteam/pkg/log"
)

// POSTReq defines the body of POST board requests.
type POSTReq struct {
	Name string `json:"name"`
}

// POSTResp defines the body of POST board responses.
type POSTResp struct {
	Error string `json:"error,omitempty"`
}

// POSTHandler is an api.MethodHandler that can be used to handle POST board
// requests.
type POSTHandler struct {
	userSelector  legacydb.Selector[userTable.Record]
	validator     api.StringValidator
	teamSelector  legacydb.Selector[teamTable.Record]
	boardCounter  legacydb.Counter
	boardInserter legacydb.Inserter[boardTable.InRecord]
	log           pkgLog.Errorer
}

// NewPOSTHandler creates and returns a new POSTHandler.
func NewPOSTHandler(
	userSelector legacydb.Selector[userTable.Record],
	validator api.StringValidator,
	boardCounter legacydb.Counter,
	boardInserter legacydb.Inserter[boardTable.InRecord],
	log pkgLog.Errorer,
) POSTHandler {
	return POSTHandler{
		userSelector:  userSelector,
		validator:     validator,
		boardCounter:  boardCounter,
		boardInserter: boardInserter,
		log:           log,
	}
}

// Handle handles the POST requests sent to the board route.
func (h POSTHandler) Handle(
	w http.ResponseWriter, r *http.Request, username string,
) {
	// Validate that the user is a team admin.
	user, err := h.userSelector.Select(username)
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusUnauthorized)
		if encodeErr := json.NewEncoder(w).Encode(
			POSTResp{Error: "Username is not recognised."},
		); encodeErr != nil {
			h.log.Error(encodeErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !user.IsAdmin {
		w.WriteHeader(http.StatusForbidden)
		if encodeErr := json.NewEncoder(w).Encode(
			POSTResp{Error: "Only team admins can create boards."},
		); encodeErr != nil {
			h.log.Error(encodeErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// Read and validate request body.
	var req POSTReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := h.validator.Validate(req.Name); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if encodeErr := json.NewEncoder(w).Encode(
			POSTResp{Error: err.Error()},
		); encodeErr != nil {
			h.log.Error(encodeErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// Validate that the user has less than 3 boards. This is done to limit the
	// resources used by this demo app.
	boardCount, err := h.boardCounter.Count(strconv.Itoa(user.TeamID))
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if boardCount >= 3 {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(
			POSTResp{
				Error: "You have already created the maximum amount of " +
					"boards allowed per user. Please delete one of your " +
					"boards to create a new one.",
			},
		); err != nil {
			h.log.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// Create a new board.
	if err := h.boardInserter.Insert(
		boardTable.NewInRecord(req.Name, user.TeamID),
	); err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// All went well. Return 200.
	w.WriteHeader(http.StatusOK)
}
