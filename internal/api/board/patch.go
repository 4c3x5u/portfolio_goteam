package board

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kxplxn/goteam/internal/api"
	"github.com/kxplxn/goteam/pkg/db"
	teamTable "github.com/kxplxn/goteam/pkg/db/team"
	pkgLog "github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/pkg/token"
)

// PatchReq defines the body of PATCH board requests.
type PatchReq teamTable.Board

// PatchResp defines the body of PATCH board responses.
type PatchResp struct {
	Error string `json:"error,omitempty"`
}

// PatchHandler can be used to handle PATCH board requests.
type PatchHandler struct {
	decodeAuth    token.DecodeFunc[token.Auth]
	decodeState   token.DecodeFunc[token.State]
	idValidator   api.StringValidator
	nameValidator api.StringValidator
	boardUpdater  db.UpdaterDualKey[teamTable.Board]
	log           pkgLog.Errorer
}

// DeleteHandler is an api.MethodHandler that can be used to handle DELETE board
// requests.
func NewPatchHandler(
	decodeAuth token.DecodeFunc[token.Auth],
	decodeState token.DecodeFunc[token.State],
	idValidator api.StringValidator,
	nameValidator api.StringValidator,
	boardUpdater db.UpdaterDualKey[teamTable.Board],
	log pkgLog.Errorer,
) *PatchHandler {
	return &PatchHandler{
		decodeAuth:    decodeAuth,
		decodeState:   decodeState,
		idValidator:   idValidator,
		nameValidator: nameValidator,
		boardUpdater:  boardUpdater,
		log:           log,
	}
}

// Handle handles PATCH board requests.
func (h *PatchHandler) Handle(
	w http.ResponseWriter, r *http.Request, username string,
) {
	// get auth token
	ckAuth, err := r.Cookie(token.AuthName)
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusUnauthorized)
		if err := json.NewEncoder(w).Encode(
			PatchResp{Error: "Auth token not found."},
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

	// decode auth token
	auth, err := h.decodeAuth(ckAuth.Value)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		if err := json.NewEncoder(w).Encode(
			PatchResp{Error: "Invalid auth token."},
		); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	}

	// validate user is admin
	if !auth.IsAdmin {
		w.WriteHeader(http.StatusForbidden)
		if err = json.NewEncoder(w).Encode(PatchResp{
			Error: "Only team admins can edit boards.",
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
			return
		}
	}

	// get state token
	ckState, err := r.Cookie(token.StateName)
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusForbidden)
		if err = json.NewEncoder(w).Encode(PatchResp{
			Error: "State token not found.",
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}

	// decode state token
	state, err := h.decodeState(ckState.Value)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		if err = json.NewEncoder(w).Encode(PatchResp{
			Error: "Invalid state token.",
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	}

	// decode board
	var req PatchReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}

	// validate board ID
	if err := h.idValidator.Validate(req.ID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		var msg string
		if errors.Is(err, ErrEmpty) {
			msg = "Board ID cannot be empty."
		} else if errors.Is(err, ErrNotUUID) {
			msg = "Board ID must be a UUID."
		}

		if err = json.NewEncoder(w).Encode(PatchResp{
			Error: msg,
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	}
	if err := h.nameValidator.Validate(req.Name); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		var msg string
		if errors.Is(err, ErrEmpty) {
			msg = "Board name cannot be empty."
		} else if errors.Is(err, ErrTooLong) {
			msg = "Board name cannot be longer than 35 characters."
		}

		if err := json.NewEncoder(w).Encode(
			PatchResp{Error: msg},
		); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	}

	// validate board access
	var hasAccess bool
	for _, b := range state.Boards {
		if b.ID == req.ID {
			hasAccess = true
			break
		}
	}
	if !hasAccess {
		w.WriteHeader(http.StatusForbidden)
		if err := json.NewEncoder(w).Encode(
			PatchResp{Error: "You do not have access to this board."},
		); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	}

	// update the board for the team
	if err := h.boardUpdater.Update(
		r.Context(), auth.TeamID, teamTable.Board(req),
	); errors.Is(err, db.ErrNoItem) {
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(
			PatchResp{Error: "Board not found."},
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
