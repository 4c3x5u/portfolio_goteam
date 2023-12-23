package tasks

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kxplxn/goteam/pkg/cookie"
	"github.com/kxplxn/goteam/pkg/db"
	"github.com/kxplxn/goteam/pkg/db/tasktable"
	"github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/pkg/validator"
)

// PatchReq defines body of PATCH tasks requests.
type PatchReq []Task

// Task represents an element in PatchReq.
type Task struct {
	ID      string    `json:"id"`
	Title   string    `json:"title"`
	Descr   string    `json:"description"`
	Order   int       `json:"order"`
	Subt    []Subtask `json:"subtasks"`
	BoardID string    `json:"board"`
	ColNo   int       `json:"column"`
}

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
	stateDecoder   cookie.Decoder[cookie.State]
	colNoValidator validator.Int
	tasksUpdater   db.Updater[[]tasktable.Task]
	stateEncoder   cookie.Encoder[cookie.State]
	log            log.Errorer
}

// NewPatchHandler creates and returns a new PATCHHandler.
func NewPatchHandler(
	authDecoder cookie.Decoder[cookie.Auth],
	stateDecoder cookie.Decoder[cookie.State],
	colNoValidator validator.Int,
	tasksUpdater db.Updater[[]tasktable.Task],
	stateEncoder cookie.Encoder[cookie.State],
	log log.Errorer,
) PatchHandler {
	return PatchHandler{
		authDecoder:    authDecoder,
		stateDecoder:   stateDecoder,
		colNoValidator: colNoValidator,
		tasksUpdater:   tasksUpdater,
		stateEncoder:   stateEncoder,
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
		if encodeErr := json.NewEncoder(w).Encode(PatchResp{
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
		if err = json.NewEncoder(w).Encode(PatchResp{
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
		if err = json.NewEncoder(w).Encode(PatchResp{
			Error: "Only team admins can edit tasks.",
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
	state, err := h.stateDecoder.Decode(*ckState)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err = json.NewEncoder(w).Encode(PatchResp{
			Error: "Invalid state token.",
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	}

	// Decode request body and map it into tasks.
	var req PatchReq
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}

	// validate task access and column numbers
	var tasks []tasktable.Task
	for _, t := range req {
		var hasAccess bool
		for _, sb := range state.Boards {
			for _, sc := range sb.Columns {
				for _, st := range sc.Tasks {
					if st.ID == t.ID {
						hasAccess = true
						break
					}
				}
				if hasAccess {
					break
				}
			}
			if hasAccess {
				break
			}
		}

		if !hasAccess {
			w.WriteHeader(http.StatusBadRequest)
			if err = json.NewEncoder(w).Encode(
				PatchResp{Error: "Invalid task ID."},
			); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				h.log.Error(err.Error())
			}
			return
		}

		if err := h.colNoValidator.Validate(t.ColNo); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			if err = json.NewEncoder(w).Encode(
				PatchResp{Error: "Invalid column number."},
			); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				h.log.Error(err.Error())
			}
			return
		}

		var subtasks []tasktable.Subtask
		for _, st := range t.Subt {
			subtasks = append(subtasks, tasktable.NewSubtask(st.Title, st.IsDone))
		}
		tasks = append(tasks, tasktable.NewTask(
			auth.TeamID,
			t.BoardID,
			t.ColNo,
			t.ID,
			t.Title,
			t.Descr,
			t.Order,
			subtasks,
		))
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
			h.log.Error(err.Error())
		}
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}

	// generate new state
	// FIXME: this is wrong but don't fix it for now because I'm going to
	//        replace state tokens with signed-IDs soon
	var boards []cookie.Board
	for _, stB := range state.Boards {
		var columns []cookie.Column
		for _, stC := range stB.Columns {
			var tasks []cookie.Task
			for _, stT := range stC.Tasks {
				tasks = append(tasks, stT)
			}
			columns = append(columns, cookie.NewColumn(tasks))
		}
		boards = append(boards, cookie.NewBoard(stB.ID, columns))
	}
	newState := cookie.NewState(boards)

	// encode state into cookie
	outCkState, err := h.stateEncoder.Encode(newState)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}
	http.SetCookie(w, &outCkState)
}
