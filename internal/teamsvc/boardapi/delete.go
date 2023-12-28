package boardapi

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/kxplxn/goteam/pkg/cookie"
	"github.com/kxplxn/goteam/pkg/db"
	"github.com/kxplxn/goteam/pkg/log"
)

// DeleteHandler is an api.MethodHandler that can be used to handle DELETE board
// requests.
type DeleteHandler struct {
	authDecoder  cookie.Decoder[cookie.Auth]
	boardDeleter db.DeleterDualKey
	log          log.Errorer
}

// NewDeleteHandler creates and returns a new DeleteHandler.
func NewDeleteHandler(
	authDecoder cookie.Decoder[cookie.Auth],
	boardDeleter db.DeleterDualKey,
	log log.Errorer,
) DeleteHandler {
	return DeleteHandler{
		authDecoder:  authDecoder,
		boardDeleter: boardDeleter,
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
		h.log.Error(err)
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

	// validate ID
	id := r.URL.Query().Get("id")
	if _, err := uuid.Parse(id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// delete the board
	if err = h.boardDeleter.Delete(
		r.Context(), auth.TeamID, id,
	); errors.Is(err, db.ErrNoItem) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		h.log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
