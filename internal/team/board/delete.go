package board

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/kxplxn/goteam/pkg/cookie"
	"github.com/kxplxn/goteam/pkg/db"
	pkgLog "github.com/kxplxn/goteam/pkg/log"
)

// DeleteHandler is an api.MethodHandler that can be used to handle DELETE board
// requests.
type DeleteHandler struct {
	authDecoder  cookie.Decoder[cookie.Auth]
	stateDecoder cookie.Decoder[cookie.State]
	boardDeleter db.DeleterDualKey
	stateEncoder cookie.Encoder[cookie.State]
	log          pkgLog.Errorer
}

// NewDeleteHandler creates and returns a new DeleteHandler.
func NewDeleteHandler(
	authDecoder cookie.Decoder[cookie.Auth],
	stateDecoder cookie.Decoder[cookie.State],
	boardDeleter db.DeleterDualKey,
	stateEncoder cookie.Encoder[cookie.State],
	log pkgLog.Errorer,
) DeleteHandler {
	return DeleteHandler{
		authDecoder:  authDecoder,
		stateDecoder: stateDecoder,
		boardDeleter: boardDeleter,
		stateEncoder: stateEncoder,
		log:          log,
	}
}

// Handle handles DELETE board requests.
func (h DeleteHandler) Handle(
	w http.ResponseWriter, r *http.Request, username string,
) {
	// get auth token
	ckAuth, err := r.Cookie(cookie.AuthName)
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}

	// decode auth token
	auth, err := h.authDecoder.Decode(*ckAuth)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// validate user is admin
	if !auth.IsAdmin {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// get state token
	ckState, err := r.Cookie(cookie.StateName)
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusForbidden)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}

	// decode state token
	state, err := h.stateDecoder.Decode(*ckState)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// get id and check it's a valid GUID
	id := r.URL.Query().Get("id")
	if _, err := uuid.Parse(id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// check if the user has access to the board
	var hasAccess bool
	var newBoards []cookie.Board
	for _, b := range state.Boards {
		if b.ID == id {
			hasAccess = true
			continue
		}
		newBoards = append(newBoards, b)
	}
	if !hasAccess {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	state.Boards = newBoards

	// delete the board
	if err = h.boardDeleter.Delete(r.Context(), auth.TeamID, id); err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// encode the new state
	outCkState, err := h.stateEncoder.Encode(state)
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &outCkState)
}
