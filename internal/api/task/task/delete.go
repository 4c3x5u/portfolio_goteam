package task

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/kxplxn/goteam/pkg/db"
	pkgLog "github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/pkg/token"
)

// DeleteResp defines the body of DELETE task responses.
type DeleteResp struct {
	Error string `json:"error"`
}

// DeleteHandler is an api.MethodHandler that can be used to handle DELETE
// requests made to the task route.
type DeleteHandler struct {
	decodeAuth  token.DecodeFunc[token.Auth]
	decodeState token.DecodeFunc[token.State]
	taskDeleter db.DeleterDualKey
	encodeState token.EncodeFunc[token.State]
	log         pkgLog.Errorer
}

// NewDeleteHandler creates and returns a new DELETEHandler.
func NewDeleteHandler(
	decodeAuth token.DecodeFunc[token.Auth],
	decodeState token.DecodeFunc[token.State],
	taskDeleter db.DeleterDualKey,
	encodeState token.EncodeFunc[token.State],
	log pkgLog.Errorer,
) DeleteHandler {
	return DeleteHandler{
		decodeAuth:  decodeAuth,
		decodeState: decodeState,
		taskDeleter: taskDeleter,
		encodeState: encodeState,
		log:         log,
	}
}

// Handle handles the DELETE requests sent to the task route.
func (h DeleteHandler) Handle(
	w http.ResponseWriter, r *http.Request, username string,
) {
	// get auth token
	ckAuth, err := r.Cookie(token.AuthName)
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusUnauthorized)
		if encodeErr := json.NewEncoder(w).Encode(DeleteResp{
			Error: "Auth token not found.",
		}); encodeErr != nil {
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
		if err = json.NewEncoder(w).Encode(DeleteResp{
			Error: "Invalid auth token.",
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
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
			h.log.Error(err.Error())
			return
		}
	}

	// get state token
	ckState, err := r.Cookie(token.StateName)
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusBadRequest)
		if encodeErr := json.NewEncoder(w).Encode(DeleteResp{
			Error: "State token not found.",
		}); encodeErr != nil {
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
		w.WriteHeader(http.StatusBadRequest)
		if err = json.NewEncoder(w).Encode(DeleteResp{
			Error: "Invalid state token.",
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
			return
		}
	}

	// validate task ID exists in state
	id := r.URL.Query().Get("id")
	var idValid bool
	for _, b := range state.Boards {
		for _, c := range b.Columns {
			for _, t := range c.Tasks {
				if t.ID == id {
					idValid = true
					break
				}
			}
			if idValid {
				break
			}
		}
		if idValid {
			break
		}
	}
	if !idValid {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(DeleteResp{
			Error: "Invalid task ID.",
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
			return
		}
	}

	// delete task from the task table
	if err = h.taskDeleter.Delete(
		r.Context(), auth.TeamID, id,
	); errors.Is(err, db.ErrNoItem) {
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(DeleteResp{
			Error: "Task not found.",
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
			return
		}
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}

	// update state
	for _, b := range state.Boards {
		for _, c := range b.Columns {
			var tasks []token.Task
			for _, t := range c.Tasks {
				if t.ID != id {
					tasks = append(tasks, t)
				}
			}
			c.Tasks = tasks
		}
	}

	// encode state into cookie
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
