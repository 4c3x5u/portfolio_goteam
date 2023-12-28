package tasksapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kxplxn/goteam/pkg/cookie"
	"github.com/kxplxn/goteam/pkg/db"
	"github.com/kxplxn/goteam/pkg/db/tasktbl"
	"github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/pkg/validator"
)

// GetResp defines the body of GET tasks responses.
type GetResp []tasktbl.Task

// GetHandler is an api.MethodHandler that can handle GET requests sent to the
// tasks route.
type GetHandler struct {
	boardIDValidator validator.String
	stateDecoder     cookie.Decoder[cookie.State]
	retrieverByBoard db.Retriever[[]tasktbl.Task]
	authDecoder      cookie.Decoder[cookie.Auth]
	retrieverByTeam  db.Retriever[[]tasktbl.Task]
	stateEncoder     cookie.Encoder[cookie.State]
	log              log.Errorer
}

// NewGetHandler creates and returns a new GetHandler.
func NewGetHandler(
	boardIDValidator validator.String,
	retrieverByBoard db.Retriever[[]tasktbl.Task],
	authDecoder cookie.Decoder[cookie.Auth],
	retrieverByTeam db.Retriever[[]tasktbl.Task],
	log log.Errorer,
) GetHandler {
	return GetHandler{
		boardIDValidator: boardIDValidator,
		retrieverByBoard: retrieverByBoard,
		authDecoder:      authDecoder,
		retrieverByTeam:  retrieverByTeam,
		log:              log,
	}
}

// Handle handles GET requests sent to the tasks route.
func (h GetHandler) Handle(w http.ResponseWriter, r *http.Request, _ string) {
	// get auth token
	ckAuth, err := r.Cookie(cookie.AuthName)
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// decode state token
	auth, err := h.authDecoder.Decode(*ckAuth)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		h.log.Error(err)
		return
	}

	// get tasks by board ID if present, otherwise get tasks by team ID of the
	// auth cookie
	var (
		tasks  []tasktbl.Task
		status int
	)
	if boardID := r.URL.Query().Get("boardID"); boardID != "" {
		tasks, status = h.getByBoardID(r.Context(), auth, w, boardID)
	} else {
		tasks, status = h.getByTeamID(r.Context(), auth, w)
	}

	// write status and if OK, write tasks to response
	w.WriteHeader(status)
	if status == http.StatusOK {
		if err := json.NewEncoder(w).Encode(tasks); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err)
			return
		}
	}
}

// getByBoardID validates the board ID and retrieves all tasks for the board,
// writing them to the response.
func (h GetHandler) getByBoardID(
	ctx context.Context, auth cookie.Auth, w http.ResponseWriter, boardID string,
) ([]tasktbl.Task, int) {
	if err := h.boardIDValidator.Validate(boardID); err != nil {
		return nil, http.StatusBadRequest
	}

	// retrieve tasks
	tasks, err := h.retrieverByBoard.Retrieve(ctx, boardID)
	if errors.Is(err, db.ErrNoItem) {
		// if no items, set tasks to empty slice
		tasks = []tasktbl.Task{}
	} else if err != nil {
		h.log.Error(err)
		return nil, http.StatusInternalServerError
	}

	// validate that all tasks belong to user's team
	for _, t := range tasks {
		if t.TeamID != auth.TeamID {
			return nil, http.StatusForbidden
		}
	}

	// return tasks
	return tasks, http.StatusOK
}

// getByTeamID gets the team ID from the auth token, retrieves all tasks for
// the team, and writes the ones with the first task's board ID to the response.
func (h GetHandler) getByTeamID(
	ctx context.Context, auth cookie.Auth, w http.ResponseWriter,
) ([]tasktbl.Task, int) {
	// retrieve tasks
	tasks, err := h.retrieverByTeam.Retrieve(ctx, auth.TeamID)
	if errors.Is(err, db.ErrNoItem) {
		// if no items, set tasks to empty slice
		tasks = []tasktbl.Task{}
	} else if err != nil {
		h.log.Error(err)
		return nil, http.StatusInternalServerError
	}

	// if more than one task, only return the ones with the first task's board
	// ID
	if len(tasks) > 1 {
		singleBoardTasks := []tasktbl.Task{}
		var boardID string
		for _, t := range tasks {
			switch boardID {
			case "":
				boardID = t.BoardID
				singleBoardTasks = append(singleBoardTasks, t)
			case t.BoardID:
				singleBoardTasks = append(singleBoardTasks, t)
			}
		}
		tasks = singleBoardTasks
	}

	return tasks, http.StatusOK
}
