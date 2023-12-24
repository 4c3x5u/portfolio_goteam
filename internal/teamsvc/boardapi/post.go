package boardapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/kxplxn/goteam/pkg/cookie"
	"github.com/kxplxn/goteam/pkg/db"
	"github.com/kxplxn/goteam/pkg/db/teamtbl"
	"github.com/kxplxn/goteam/pkg/log"
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
	authDecoder   cookie.Decoder[cookie.Auth]
	stateDecoder  cookie.Decoder[cookie.State]
	nameValidator validator.String
	inserter      db.InserterDualKey[teamtbl.Board]
	stateEncoder  cookie.Encoder[cookie.State]
	log           log.Errorer
}

// NewPostHandler creates and returns a new PostHandler.
func NewPostHandler(
	authDecoder cookie.Decoder[cookie.Auth],
	stateDecoder cookie.Decoder[cookie.State],
	nameValidator validator.String,
	inserter db.InserterDualKey[teamtbl.Board],
	stateEncoder cookie.Encoder[cookie.State],
	log log.Errorer,
) *PostHandler {
	return &PostHandler{
		authDecoder:   authDecoder,
		stateDecoder:  stateDecoder,
		nameValidator: nameValidator,
		inserter:      inserter,
		stateEncoder:  stateEncoder,
		log:           log,
	}
}

// Handle handles DELETE board requests.
func (h PostHandler) Handle(
	w http.ResponseWriter, r *http.Request, username string,
) {
	// get auth token
	ckAuth, err := r.Cookie(cookie.AuthName)
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusUnauthorized)
		if err := json.NewEncoder(w).Encode(
			PatchResp{Error: "Auth token not found."},
		); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err)
		}
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err)
		return
	}

	// decode auth token
	auth, err := h.authDecoder.Decode(*ckAuth)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		if err := json.NewEncoder(w).Encode(
			PatchResp{Error: "Invalid auth token."},
		); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err)
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
			h.log.Error(err)
		}
		return
	}

	// get state token
	ckState, err := r.Cookie(cookie.StateName)
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusForbidden)
		if err = json.NewEncoder(w).Encode(PatchResp{
			Error: "State token not found.",
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err)
		}
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err)
		return
	}

	// decode state token
	state, err := h.stateDecoder.Decode(*ckState)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		if err = json.NewEncoder(w).Encode(PatchResp{
			Error: "Invalid state token.",
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err)
		}
		return
	}

	// check if the user's team already has 3 boards
	if len(state.Boards) > 2 {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(
			PostResp{Error: msgLimitReached},
		); err != nil {
			h.log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// get and validate board name
	var req PostReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error(err)
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
			h.log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// insert the board into the team's boards in the team table - retry up to 3
	// times for the unlikely event that the generated UUID is a duplicate
	id := uuid.NewString()
	for i := 0; i < 3; i++ {
		if err = h.inserter.Insert(r.Context(), auth.TeamID, teamtbl.Board{
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
				h.log.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		} else if err != nil {
			h.log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// update, encode, and set state token
	state.Boards = append(state.Boards, cookie.Board{
		ID: id, Columns: []cookie.Column{
			{Tasks: []cookie.Task{}},
			{Tasks: []cookie.Task{}},
			{Tasks: []cookie.Task{}},
			{Tasks: []cookie.Task{}},
		},
	})
	outCkState, err := h.stateEncoder.Encode(state)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err)
		return
	}
	http.SetCookie(w, &outCkState)

}

// msgLimitReached is the error message written into PostResp when the user's
// team already has 3 boards.
const msgLimitReached = "You have already created the maximum amount of " +
	"boards allowed per team. Please delete one of your boards to create a " +
	"new one."
