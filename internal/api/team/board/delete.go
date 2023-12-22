package board

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/kxplxn/goteam/pkg/db"
	pkgLog "github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/pkg/token"
)

// DeleteHandler is an api.MethodHandler that can be used to handle DELETE board
// requests.
type DeleteHandler struct {
	decodeAuth   token.DecodeFunc[token.Auth]
	decodeState  token.DecodeFunc[token.State]
	boardDeleter db.DeleterDualKey
	log          pkgLog.Errorer
}

// NewDeleteHandler creates and returns a new DeleteHandler.
func NewDeleteHandler(
	decodeAuth token.DecodeFunc[token.Auth],
	decodeState token.DecodeFunc[token.State],
	boardDeleter db.DeleterDualKey,
	log pkgLog.Errorer,
) DeleteHandler {
	return DeleteHandler{
		decodeAuth:   decodeAuth,
		decodeState:  decodeState,
		boardDeleter: boardDeleter,
		log:          log,
	}
}

// Handle handles DELETE board requests.
func (h DeleteHandler) Handle(
	w http.ResponseWriter, r *http.Request, username string,
) {
	// get auth token
	ckAuth, err := r.Cookie(token.AuthName)
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusUnauthorized)
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
		return
	}

	// validate user is admin
	if !auth.IsAdmin {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// get state token
	ckState, err := r.Cookie(token.StateName)
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusForbidden)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}

	// decode auth token
	state, err := h.decodeState(ckState.Value)
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
	for _, b := range state.Boards {
		if b.ID == id {
			hasAccess = true
			break
		}
	}
	if !hasAccess {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// delete the board
	if err = h.boardDeleter.Delete(r.Context(), auth.TeamID, id); err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
