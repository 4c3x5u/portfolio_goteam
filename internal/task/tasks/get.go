package tasks

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kxplxn/goteam/pkg/cookie"
	"github.com/kxplxn/goteam/pkg/db"
	"github.com/kxplxn/goteam/pkg/db/tasktable"
	"github.com/kxplxn/goteam/pkg/log"
)

// GetResp defines the body of GET tasks responses.
type GetResp []tasktable.Task

// GetHandler is an api.MethodHandler that can handle GET requests sent to the
// tasks route.
type GetHandler struct {
	authDecoder cookie.Decoder[cookie.Auth]
	retriever   db.Retriever[[]tasktable.Task]
	log         log.Errorer
}

// NewGetHandler creates and returns a new GetHandler.
func NewGetHandler(
	authDecoder cookie.Decoder[cookie.Auth],
	retriever db.Retriever[[]tasktable.Task],
	log log.Errorer,
) GetHandler {
	return GetHandler{authDecoder: authDecoder, retriever: retriever, log: log}
}

// Handle handles GET requests sent to the tasks route.
func (h GetHandler) Handle(w http.ResponseWriter, r *http.Request, _ string) {
	// get auth token
	ckAuth, err := r.Cookie(cookie.AuthName)
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// decode auth token
	auth, err := h.authDecoder.Decode(*ckAuth)
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
