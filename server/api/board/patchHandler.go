package board

import (
	"encoding/json"
	"net/http"
	pkgLog "server/log"
)

type PATCHHandler struct {
	idValidator StringValidator
	log         pkgLog.Errorer
}

func NewPATCHHandler(
	idValidator StringValidator, log pkgLog.Errorer,
) *PATCHHandler {
	return &PATCHHandler{idValidator: idValidator, log: log}
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
}
