package taskapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kxplxn/goteam/pkg/cookie"
	"github.com/kxplxn/goteam/pkg/db"
	"github.com/kxplxn/goteam/pkg/log"
)

// DeleteResp defines the body of DELETE task responses.
type DeleteResp struct {
	Error string `json:"error"`
}

// DeleteHandler is an api.MethodHandler that can be used to handle DELETE
// requests made to the task route.
type DeleteHandler struct {
	authDecoder cookie.Decoder[cookie.Auth]
	taskDeleter db.DeleterDualKey
	log         log.Errorer
}

// NewDeleteHandler creates and returns a new DELETEHandler.
func NewDeleteHandler(
	authDecoder cookie.Decoder[cookie.Auth],
	taskDeleter db.DeleterDualKey,
	log log.Errorer,
) DeleteHandler {
	return DeleteHandler{
		authDecoder: authDecoder,
		taskDeleter: taskDeleter,
		log:         log,
	}
}

// Handle handles the DELETE requests sent to the task route.
func (h DeleteHandler) Handle(
	w http.ResponseWriter, r *http.Request, username string,
) {
	// get auth token
	ckAuth, err := r.Cookie(cookie.AuthName)
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusUnauthorized)
		if encodeErr := json.NewEncoder(w).Encode(DeleteResp{
			Error: "Auth token not found.",
		}); encodeErr != nil {
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
		if err = json.NewEncoder(w).Encode(DeleteResp{
			Error: "Invalid auth token.",
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err)
			return
		}
	}

	// validate user is admin
	if !auth.IsAdmin {
		w.WriteHeader(http.StatusForbidden)
		if err = json.NewEncoder(w).Encode(DeleteResp{
			Error: "Only team admins can delete tasks.",
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err)
			return
		}
	}

	// delete task from the task table
	if err = h.taskDeleter.Delete(
		r.Context(), auth.TeamID, r.URL.Query().Get("id"),
	); errors.Is(err, db.ErrNoItem) {
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(DeleteResp{
			Error: "Task not found.",
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err)
			return
		}
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err)
		return
	}
}
