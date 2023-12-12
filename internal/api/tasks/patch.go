package tasks

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kxplxn/goteam/internal/api"
	"github.com/kxplxn/goteam/pkg/dbaccess"
	columnTable "github.com/kxplxn/goteam/pkg/dbaccess/column"
	pkgLog "github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/pkg/token"
)

// PatchReq defines body of PATCH tasks requests.
type PatchReq []Task

// Task represents an element in PatchReq.
type Task struct {
	ID    string `json:"id"`
	Order int    `json:"order"`
	ColNo int    `json:"columnNumber"`
}

// PatchResp defines the body for PATCH column responses.
type PatchResp struct {
	Error string `json:"error"`
}

// PatchHandler is an api.MethodHandler that can be used to handle PATCH
// requests sent to the tasks route.
type PatchHandler struct {
	decodeAuth     token.DecodeFunc[token.Auth]
	decodeState    token.DecodeFunc[token.State]
	colNoValidator api.IntValidator
	columnUpdater  dbaccess.Updater[[]columnTable.Task]
	log            pkgLog.Errorer
}

// NewPATCHHandler creates and returns a new PATCHHandler.
func NewPATCHHandler(
	decodeAuth token.DecodeFunc[token.Auth],
	decodeState token.DecodeFunc[token.State],
	colNoValidator api.IntValidator,
	columnUpdater dbaccess.Updater[[]columnTable.Task],
	log pkgLog.Errorer,
) PatchHandler {
	return PatchHandler{
		decodeAuth:     decodeAuth,
		decodeState:    decodeState,
		colNoValidator: colNoValidator,
		columnUpdater:  columnUpdater,
		log:            log,
	}
}

// Handle handles the PATCH requests sent to the tasks route.
func (h PatchHandler) Handle(
	w http.ResponseWriter, r *http.Request, username string,
) {
	// get auth token
	ckAuth, err := r.Cookie(token.AuthName)
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
	auth, err := h.decodeAuth(ckAuth.Value)
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
	ckState, err := r.Cookie(token.StateName)
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
	state, err := h.decodeState(ckState.Value)
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
	// TODO: remove
	var tasks []columnTable.Task
	for _, t := range req {
		tasks = append(tasks, columnTable.Task{ID: 0, Order: t.Order})
	}

	// validate task access and column numbers
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

	}

	// update task records in the database using column ID and order from tasks
	if err = h.columnUpdater.Update("", tasks); errors.Is(err, sql.ErrNoRows) {
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
}
