package tasks

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kxplxn/goteam/pkg/db"
	"github.com/kxplxn/goteam/pkg/db/tasktable"
	pkgLog "github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/pkg/token"
)

// GetResp defines the body of GET tasks responses.
type GetResp []tasktable.Task

// GetHandler is an api.MethodHandler that can handle GET requests sent to the
// tasks route.
type GetHandler struct {
	decodeAuth token.DecodeFunc[token.Auth]
	retriever  db.Retriever[[]tasktable.Task]
	log        pkgLog.Errorer
}

// NewGetHandler creates and returns a new GetHandler.
func NewGetHandler(
	decodeAuth token.DecodeFunc[token.Auth],
	retriever db.Retriever[[]tasktable.Task],
	log pkgLog.Errorer,
) GetHandler {
	return GetHandler{decodeAuth: decodeAuth, retriever: retriever, log: log}
}

// Handle handles GET requests sent to the tasks route.
func (h GetHandler) Handle(w http.ResponseWriter, r *http.Request, _ string) {
	// get auth token
	ckAuth, err := r.Cookie(token.AuthName)
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// decode auth token
	auth, err := h.decodeAuth(ckAuth.Value)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// retrieve tasks
	tasks, err := h.retriever.Retrieve(r.Context(), auth.TeamID)
	if errors.Is(err, db.ErrNoItem) {
		// if no items, set tasks to empty slice
		tasks = []tasktable.Task{}
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// write response
	if err = json.NewEncoder(w).Encode(tasks); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}
}
