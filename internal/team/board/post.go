package board

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/kxplxn/goteam/pkg/db"
	"github.com/kxplxn/goteam/pkg/db/teamtable"
	pkgLog "github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/pkg/token"
	"github.com/kxplxn/goteam/pkg/validator"
)

// PostReq defines the body of POST board requests.
type PostReq struct {
	Name string `json:"name"`
}

// PostResp defines the body of POST board responses.
type PostResp struct {
	Error string `json:"error,omitempty"`
}

// DeleteHandler is an api.MethodHandler that can be used to handle POST board
// requests.
type PostHandler struct {
	decodeAuth    token.DecodeFunc[token.Auth]
	decodeState   token.DecodeFunc[token.State]
	nameValidator validator.String
	inserter      db.InserterDualKey[teamtable.Board]
	encodeState   token.EncodeFunc[token.State]
	log           pkgLog.Errorer
}

// NewPostHandler creates and returns a new PostHandler.
func NewPostHandler(
	decodeAuth token.DecodeFunc[token.Auth],
	decodeState token.DecodeFunc[token.State],
	nameValidator validator.String,
	inserter db.InserterDualKey[teamtable.Board],
	encodeState token.EncodeFunc[token.State],
	log pkgLog.Errorer,
) *PostHandler {
	return &PostHandler{
		decodeAuth:    decodeAuth,
		decodeState:   decodeState,
		nameValidator: nameValidator,
		inserter:      inserter,
		encodeState:   encodeState,
		log:           log,
	}
}

// Handle handles DELETE board requests.
func (h PostHandler) Handle(
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
		}
		return
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

	// check if the user's team already has 3 boards
	if len(state.Boards) > 2 {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(
			PostResp{Error: msgLimitReached},
		); err != nil {
			h.log.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// get and validate board name
	var req PostReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := h.nameValidator.Validate(req.Name); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		var msg string
		if errors.Is(err, validator.ErrEmpty) {
			msg = "Board name cannot be empty."
		} else if errors.Is(err, validator.ErrTooLong) {
			msg = "Board name cannot be longer than 35 characters."
		}

		if err = json.NewEncoder(w).Encode(PostResp{Error: msg}); err != nil {
			h.log.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// insert the board into the team's boards in the team table - retry up to 3
	// times for the unlikely event that the generated UUID is a duplicate
	id := uuid.NewString()
	for i := 0; i < 3; i++ {
		if err = h.inserter.Insert(r.Context(), auth.TeamID, teamtable.Board{
			ID:   id,
			Name: req.Name,
		}); errors.Is(err, db.ErrDupKey) {
			id = uuid.NewString()
			continue
		} else if errors.Is(err, db.ErrLimitReached) {
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(
				PostResp{Error: msgLimitReached},
			); err != nil {
				h.log.Error(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		} else if err != nil {
			h.log.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// update, encode, and set state token
	state.Boards = append(state.Boards, token.Board{
		ID: id, Columns: []token.Column{
			{Tasks: []token.Task{}},
			{Tasks: []token.Task{}},
			{Tasks: []token.Task{}},
			{Tasks: []token.Task{}},
		},
	})
	exp := time.Now().Add(token.DefaultDuration).UTC()
	tkState, err := h.encodeState(exp, state)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     token.StateName,
		Value:    tkState,
		Expires:  exp,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})

}

// msgLimitReached is the error message written into PostResp when the user's
// team already has 3 boards.
const msgLimitReached = "You have already created the maximum amount of " +
	"boards allowed per team. Please delete one of your boards to create a " +
	"new one."
