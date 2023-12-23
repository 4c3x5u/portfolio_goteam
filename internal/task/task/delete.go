package task

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
	authDecoder  cookie.Decoder[cookie.Auth]
	stateDecoder cookie.Decoder[cookie.State]
	taskDeleter  db.DeleterDualKey
	stateEncoder cookie.Encoder[cookie.State]
	log          log.Errorer
}

// NewDeleteHandler creates and returns a new DELETEHandler.
func NewDeleteHandler(
	authDecoder cookie.Decoder[cookie.Auth],
	stateDecoder cookie.Decoder[cookie.State],
	taskDeleter db.DeleterDualKey,
	stateEncoder cookie.Encoder[cookie.State],
	log log.Errorer,
) DeleteHandler {
	return DeleteHandler{
		authDecoder:  authDecoder,
		stateDecoder: stateDecoder,
		taskDeleter:  taskDeleter,
		stateEncoder: stateEncoder,
		log:          log,
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
			h.log.Error(err.Error())
		}
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
	ckState, err := r.Cookie(cookie.StateName)
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
	state, err := h.stateDecoder.Decode(*ckState)
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
			var tasks []cookie.Task
			for _, t := range c.Tasks {
				if t.ID != id {
					tasks = append(tasks, t)
				}
			}
			c.Tasks = tasks
		}
	}

	// encode state into cookie
	outCkState, err := h.stateEncoder.Encode(state)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}

	// set state cookie
	http.SetCookie(w, &outCkState)
}
