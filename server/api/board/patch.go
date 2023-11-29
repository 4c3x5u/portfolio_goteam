package board

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kxplxn/goteam/server/api"
	"github.com/kxplxn/goteam/server/dbaccess"
	boardTable "github.com/kxplxn/goteam/server/dbaccess/board"
	userTable "github.com/kxplxn/goteam/server/dbaccess/user"
	pkgLog "github.com/kxplxn/goteam/server/log"
)

// PATCHReq defines the request body for PATCH board requests.
type PATCHReq struct {
	Name string `json:"name"`
}

// PATCHResp defines the response body for PATCH board requests.
type PATCHResp struct {
	Error string `json:"error,omitempty"`
}

type PATCHHandler struct {
	userSelector  dbaccess.Selector[userTable.Record]
	idValidator   api.StringValidator
	nameValidator api.StringValidator
	boardSelector dbaccess.Selector[boardTable.Record]
	boardUpdater  dbaccess.Updater[string]
	log           pkgLog.Errorer
}

func NewPATCHHandler(
	userSelector dbaccess.Selector[userTable.Record],
	idValidator api.StringValidator,
	nameValidator api.StringValidator,
	boardSelector dbaccess.Selector[boardTable.Record],
	boardUpdater dbaccess.Updater[string],
	log pkgLog.Errorer,
) *PATCHHandler {
	return &PATCHHandler{
		userSelector:  userSelector,
		idValidator:   idValidator,
		nameValidator: nameValidator,
		boardSelector: boardSelector,
		boardUpdater:  boardUpdater,
		log:           log,
	}
}

func (h *PATCHHandler) Handle(
	w http.ResponseWriter, r *http.Request, username string,
) {
	// Validate that the user is a team admin.
	user, err := h.userSelector.Select(username)
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusUnauthorized)
		if err := json.NewEncoder(w).Encode(
			PATCHResp{Error: "Username is not recognised."},
		); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}
	if !user.IsAdmin {
		w.WriteHeader(http.StatusForbidden)
		if err := json.NewEncoder(w).Encode(
			PATCHResp{Error: "Only team admins can edit the board."},
		); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	}

	// Retrieve and validate the board ID.
	boardID := r.URL.Query().Get("id")
	if err := h.idValidator.Validate(boardID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(
			PATCHResp{Error: err.Error()},
		); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	}

	// Retrieve and validate the new board name.
	var req PATCHReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}
	if err := h.nameValidator.Validate(req.Name); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(
			PATCHResp{Error: err.Error()},
		); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	}

	// Validate that the board exists in the database and that it belongs to the
	// team that the user is the admin of.
	board, err := h.boardSelector.Select(boardID)
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(
			PATCHResp{Error: "Board not found."},
		); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}
	if board.TeamID != user.TeamID {
		w.WriteHeader(http.StatusForbidden)
		if err := json.NewEncoder(w).Encode(
			PATCHResp{Error: "You do not have access to this board."},
		); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	}

	// Update the name of the board in the database. Note that 500 is returned
	// any error including sql.ErrNoRows. Because even in that case,
	// something else must have gone wrong since we have already validated that
	// a board with this ID exists.
	if err := h.boardUpdater.Update(boardID, req.Name); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}

	// All went well. Return 200.
	w.WriteHeader(http.StatusOK)
}
