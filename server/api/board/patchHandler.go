package board

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"server/dbaccess"
	pkgLog "server/log"
)

type PATCHHandler struct {
	idValidator       StringValidator
	nameValidator     StringValidator
	boardSelector     dbaccess.Selector[dbaccess.Board]
	userBoardSelector dbaccess.RelSelector[bool]
	boardUpdater      dbaccess.Updater
	log               pkgLog.Errorer
}

func NewPATCHHandler(
	idValidator StringValidator,
	nameValidator StringValidator,
	boardSelector dbaccess.Selector[dbaccess.Board],
	userBoardSelector dbaccess.RelSelector[bool],
	boardUpdater dbaccess.Updater,
	log pkgLog.Errorer,
) *PATCHHandler {
	return &PATCHHandler{
		idValidator:       idValidator,
		nameValidator:     nameValidator,
		boardSelector:     boardSelector,
		userBoardSelector: userBoardSelector,
		boardUpdater:      boardUpdater,
		log:               log,
	}
}

func (h *PATCHHandler) Handle(
	w http.ResponseWriter, r *http.Request, username string,
) {
	// Retrieve and validate the board ID.
	boardID := r.URL.Query().Get("id")
	if err := h.idValidator.Validate(boardID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(
			ResBody{Error: err.Error()},
		); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	}

	// Retrieve and validate the new board name.
	var reqBody ReqBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}
	if err := h.nameValidator.Validate(reqBody.Name); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(
			ResBody{Error: err.Error()},
		); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	}

	// Validate that the board exists in the database.
	if _, err := h.boardSelector.Select(
		boardID,
	); errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(
			ResBody{Error: "Board not found."},
		); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}

	// Validate that the user is a board admin.
	if isAdmin, err := h.userBoardSelector.Select(
		username, boardID,
	); errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusForbidden)
		if err := json.NewEncoder(w).Encode(
			ResBody{Error: "You do not have access to this board."},
		); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	} else if !isAdmin {
		w.WriteHeader(http.StatusForbidden)
		if err := json.NewEncoder(w).Encode(
			ResBody{Error: "Only board admins can edit the board."},
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
	if err := h.boardUpdater.Update(boardID, reqBody.Name); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}

	// All went well. Return 200.
	w.WriteHeader(http.StatusOK)
}
