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
	idValidator   StringValidator
	nameValidator StringValidator
	boardSelector dbaccess.Selector[dbaccess.Board]
	log           pkgLog.Errorer
}

func NewPATCHHandler(
	idValidator StringValidator,
	nameValidator StringValidator,
	boardSelector dbaccess.Selector[dbaccess.Board],
	log pkgLog.Errorer,
) *PATCHHandler {
	return &PATCHHandler{
		idValidator:   idValidator,
		nameValidator: nameValidator,
		boardSelector: boardSelector,
		log:           log,
	}
}

func (h *PATCHHandler) Handle(
	w http.ResponseWriter, r *http.Request, _ string,
) {
	id := r.URL.Query().Get("id")
	if err := h.idValidator.Validate(id); err != nil {
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

	if _, err := h.boardSelector.Select(id); errors.Is(err, sql.ErrNoRows) {
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
}
