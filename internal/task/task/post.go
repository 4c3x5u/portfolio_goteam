package task

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/kxplxn/goteam/pkg/cookie"
	"github.com/kxplxn/goteam/pkg/db"
	"github.com/kxplxn/goteam/pkg/db/tasktable"
	"github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/pkg/validator"
)

// PostReq defines the body of POST task requests.
type PostReq struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Subtasks     []string `json:"subtasks"`
	BoardID      string   `json:"board"`
	ColumnNumber int      `json:"column"`
}

// PostResp defines the body of POST task responses.
type PostResp struct {
	Error string `json:"error"`
}

// PostHandler is an api.MethodHandler that can be used to handle POST requests
// sent to the task route.
type PostHandler struct {
	authDecoder        cookie.Decoder[cookie.Auth]
	stateDecoder       cookie.Decoder[cookie.State]
	titleValidator     validator.String
	subtTitleValidator validator.String
	colNoValidator     validator.Int
	taskInserter       db.Inserter[tasktable.Task]
	stateEncoder       cookie.Encoder[cookie.State]
	log                log.Errorer
}

// NewPostHandler creates and returns a new POSTHandler.
func NewPostHandler(
	authDecoder cookie.Decoder[cookie.Auth],
	stateDecoder cookie.Decoder[cookie.State],
	titleValidator validator.String,
	subtTitleValidator validator.String,
	colNoValidator validator.Int,
	taskInserter db.Inserter[tasktable.Task],
	stateEncoder cookie.Encoder[cookie.State],
	log log.Errorer,
) *PostHandler {
	return &PostHandler{
		authDecoder:        authDecoder,
		stateDecoder:       stateDecoder,
		titleValidator:     titleValidator,
		subtTitleValidator: subtTitleValidator,
		colNoValidator:     colNoValidator,
		taskInserter:       taskInserter,
		stateEncoder:       stateEncoder,
		log:                log,
	}
}

// Handle handles the POST requests sent to the task route.
func (h *PostHandler) Handle(
	w http.ResponseWriter, r *http.Request, _ string,
) {
	// get auth token
	ckAuth, err := r.Cookie(cookie.AuthName)
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusUnauthorized)
		if encodeErr := json.NewEncoder(w).Encode(PostResp{
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
		if encodeErr := json.NewEncoder(w).Encode(PostResp{
			Error: "Invalid auth token.",
		}); encodeErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err)
		}
		return
	}

	// validate user is admin
	if !auth.IsAdmin {
		w.WriteHeader(http.StatusForbidden)
		if encodeErr := json.NewEncoder(w).Encode(PostResp{
			Error: "Only team admins can create tasks.",
		}); encodeErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err)
		}
		return
	}

	// get state token
	ckState, err := r.Cookie(cookie.StateName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if encodeErr := json.NewEncoder(w).Encode(PostResp{
			Error: "State token not found.",
		}); encodeErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err)
		}
		return
	}

	// decode state token
	state, err := h.stateDecoder.Decode(*ckState)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if encodeErr := json.NewEncoder(w).Encode(PostResp{
			Error: "Invalid state token.",
		}); encodeErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err)
		}
		return
	}

	// decode request
	var req PostReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err)
		return
	}

	// validate column ID
	if err := h.colNoValidator.Validate(req.ColumnNumber); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if encodeErr := json.NewEncoder(w).Encode(PostResp{
			Error: "Column number out of bounds.",
		}); encodeErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err)
		}
		return
	}

	// validate board access and determine order for the task
	var hasBoardAccess bool
	var highestOrder int
	for _, b := range state.Boards {
		if req.BoardID == b.ID {
			hasBoardAccess = true
			for _, t := range b.Columns[req.ColumnNumber].Tasks {
				if t.Order > highestOrder {
					highestOrder = t.Order
				}
			}
		}
	}
	if !hasBoardAccess {
		w.WriteHeader(http.StatusForbidden)
		if encodeErr := json.NewEncoder(w).Encode(PostResp{
			Error: "You do not have access to this board.",
		}); encodeErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err)
		}
		return
	}
	order := highestOrder + 1

	// validate task
	if err := h.titleValidator.Validate(req.Title); err != nil {
		var errMsg string
		if errors.Is(err, validator.ErrEmpty) {
			errMsg = "Task title cannot be empty."
		} else if errors.Is(err, validator.ErrTooLong) {
			errMsg = "Task title cannot be longer than 50 characters."
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		if err = json.NewEncoder(w).Encode(PostResp{
			Error: errMsg,
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err)
		}
		return
	}

	// validate subtasks
	var subtasks []tasktable.Subtask
	for _, title := range req.Subtasks {
		if err := h.subtTitleValidator.Validate(title); err != nil {
			var errMsg string
			if errors.Is(err, validator.ErrEmpty) {
				errMsg = "Subtask title cannot be empty."
			} else if errors.Is(err, validator.ErrTooLong) {
				errMsg = "Subtask title cannot be longer than 50 characters."
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				h.log.Error(err)
				return
			}

			w.WriteHeader(http.StatusBadRequest)
			if err = json.NewEncoder(w).Encode(PostResp{
				Error: errMsg,
			}); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				h.log.Error(err)
			}
			return
		}
		subtasks = append(subtasks, tasktable.Subtask{
			Title: title, IsDone: false,
		})
	}

	// insert a new task into the task table - retry up to 3 times for the
	// unlikely event that the generated UUID is a duplicate
	id := uuid.NewString()
	for tries := 0; tries < 3; tries++ {
		if err = h.taskInserter.Insert(r.Context(), tasktable.NewTask(
			auth.TeamID,
			req.BoardID,
			req.ColumnNumber,
			id,
			req.Title,
			req.Description,
			order,
			subtasks,
		)); errors.Is(err, db.ErrDupKey) {
			id = uuid.NewString()
		} else if err != nil {
			break
		}
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err)
		return
	}

	// update, encode, and set state token
	for _, b := range state.Boards {
		if b.ID == req.BoardID {
			b.Columns[req.ColumnNumber].Tasks = append(
				b.Columns[req.ColumnNumber].Tasks,
				cookie.NewTask(id, order),
			)
		}
	}
	outCkState, err := h.stateEncoder.Encode(state)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err)
		return
	}
	http.SetCookie(w, &outCkState)
}
