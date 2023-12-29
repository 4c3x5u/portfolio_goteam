package tasksapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kxplxn/goteam/pkg/cookie"
	"github.com/kxplxn/goteam/pkg/db"
	"github.com/kxplxn/goteam/pkg/db/tasktbl"
	"github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/pkg/validator"
)

// PatchReq defines body of PATCH tasks requests.
type PatchReq []tasktbl.Task

// Subtask represents a subtask element in Task.
type Subtask struct {
	Title  string `json:"title"`
	IsDone bool   `json:"done"`
}

// PatchResp defines the body for PATCH column responses.
type PatchResp struct {
	Error string `json:"error"`
}

// PatchHandler is an api.MethodHandler that can be used to handle PATCH
// requests sent to the tasks route.
type PatchHandler struct {
	authDecoder    cookie.Decoder[cookie.Auth]
	colNoValidator validator.Int
	tasksUpdater   db.Updater[[]tasktbl.Task]
	log            log.Errorer
}

// NewPatchHandler creates and returns a new PATCHHandler.
func NewPatchHandler(
	authDecoder cookie.Decoder[cookie.Auth],
	colNoValidator validator.Int,
	tasksUpdater db.Updater[[]tasktbl.Task],
	log log.Errorer,
) PatchHandler {
	return PatchHandler{
		authDecoder:    authDecoder,
		colNoValidator: colNoValidator,
		tasksUpdater:   tasksUpdater,
		log:            log,
	}
}

// Handle handles the PATCH requests sent to the tasks route.
func (h PatchHandler) Handle(
	w http.ResponseWriter, r *http.Request, username string,
) {
	// get auth token
	ckAuth, err := r.Cookie(cookie.AuthName)
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusUnauthorized)
		if err = json.NewEncoder(w).Encode(PatchResp{
			Error: "Auth token not found.",
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

	// decode auth token
	auth, err := h.authDecoder.Decode(*ckAuth)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		if err = json.NewEncoder(w).Encode(PatchResp{
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
		if err = json.NewEncoder(w).Encode(PatchResp{
			Error: "Only team admins can edit tasks.",
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err)
			return
		}
	}

	// decode request body
	var req PatchReq
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err)
		return
	}

	if len(req) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		if err = json.NewEncoder(w).Encode(PatchResp{
			Error: "No tasks provided.",
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err)
		}
		return
	}

	// map request body into tasks, validating them as we go
	var tasks []tasktbl.Task
	for _, t := range req {
		// TODO: validate other fields, too
		if err := h.colNoValidator.Validate(t.ColNo); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			if err = json.NewEncoder(w).Encode(PatchResp{
				Error: "Invalid column number.",
			}); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				h.log.Error(err)
			}
			return
		}

		task := tasktbl.Task{
			TeamID:      auth.TeamID,
			BoardID:     t.BoardID,
			ColNo:       t.ColNo,
			ID:          t.ID,
			Title:       t.Title,
			Description: t.Description,
			Order:       t.Order,
			Subtasks:    t.Subtasks,
		}

		tasks = append(tasks, task)
	}

	// update tasks in the task table
	if err = h.tasksUpdater.Update(
		r.Context(), tasks,
	); errors.Is(err, db.ErrNoItem) {
		w.WriteHeader(http.StatusNotFound)
		if err = json.NewEncoder(w).Encode(
			PatchResp{Error: "Task not found."},
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
}
