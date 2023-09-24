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
	log               pkgLog.Errorer
}

func NewPATCHHandler(
	idValidator StringValidator,
	nameValidator StringValidator,
	boardSelector dbaccess.Selector[dbaccess.Board],
	userBoardSelector dbaccess.RelSelector[bool],
	log pkgLog.Errorer,
) *PATCHHandler {
	return &PATCHHandler{
		idValidator:       idValidator,
		nameValidator:     nameValidator,
		boardSelector:     boardSelector,
		userBoardSelector: userBoardSelector,
		log:               log,
	}
}

func (h *PATCHHandler) Handle(
	w http.ResponseWriter, r *http.Request, username string,
) {
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

	if _, err := h.userBoardSelector.Select(
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
	}
}
