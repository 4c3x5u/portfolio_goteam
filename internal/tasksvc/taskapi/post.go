package taskapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/kxplxn/goteam/pkg/cookie"
	"github.com/kxplxn/goteam/pkg/db"
	"github.com/kxplxn/goteam/pkg/db/tasktbl"
	"github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/pkg/validator"
)

// PostReq defines the body of POST task requests.
type PostReq struct {
	BoardID     string `json:"board"`
	ColNo       int    `json:"column"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Subtasks    []tasktbl.Subtask
	Order       int `json:"order"`
}

// PostResp defines the body of POST task responses.
type PostResp struct {
	Error string `json:"error"`
}

// PostHandler is an api.MethodHandler that can be used to handle POST requests
// sent to the task route.
type PostHandler struct {
	authDecoder  cookie.Decoder[cookie.Auth]
	validateReq  validator.Func[PostReq]
	taskInserter db.Inserter[tasktbl.Task]
	log          log.Errorer
}

// NewPostHandler creates and returns a new POSTHandler.
func NewPostHandler(
	authDecoder cookie.Decoder[cookie.Auth],
	validateReq validator.Func[PostReq],
	taskInserter db.Inserter[tasktbl.Task],
	log log.Errorer,
) *PostHandler {
	return &PostHandler{
		authDecoder:  authDecoder,
		validateReq:  validateReq,
		taskInserter: taskInserter,
		log:          log,
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
		if err = json.NewEncoder(w).Encode(PostResp{
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
		if err = json.NewEncoder(w).Encode(PostResp{
			Error: "Invalid auth token.",
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err)
		}
		return
	}

	// validate user is admin
	if !auth.IsAdmin {
		w.WriteHeader(http.StatusForbidden)
		if err := json.NewEncoder(w).Encode(PostResp{
			Error: "Only team admins can create tasks.",
		}); err != nil {
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

	// validate request
	if err := h.validateReq(req); err != nil {
		var msg string
		switch {
		case errors.Is(err, errBoardIDEmpty):
			msg = "Board ID cannot be empty."
		case errors.Is(err, errParseBoardID):
			msg = "Board ID is must be a valid UUID."
		case errors.Is(err, errColNoOutOfBounds):
			msg = "Column number must be between 1 and 4."
		case errors.Is(err, errTitleEmpty):
			msg = "Task title cannot be empty."
		case errors.Is(err, errTitleTooLong):
			msg = "Task title cannot be longer than 50 characters."
		case errors.Is(err, errDescTooLong):
			msg = "Task description cannot be longer than 500 characters."
		case errors.Is(err, errSubtaskTitleEmpty):
			msg = "Subtask title cannot be empty."
		case errors.Is(err, errSubtaskTitleTooLong):
			msg = "Subtask title cannot be longer than 50 characters."
		case errors.Is(err, errOrderNegative):
			msg = "Order cannot be negative."
		default:
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		if err = json.NewEncoder(w).Encode(PostResp{Error: msg}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err)
		}
		return
	}

	// insert a new task into the task table - retry up to 3 times for the
	// unlikely event that the generated UUID is a duplicate
	for tries := 0; tries < 3; tries++ {
		id := uuid.NewString()
		if err = h.taskInserter.Insert(r.Context(), tasktbl.NewTask(
			auth.TeamID,
			req.BoardID,
			req.ColNo,
			id,
			req.Title,
			req.Description,
			req.Order,
			req.Subtasks,
		)); !errors.Is(err, db.ErrDupKey) {
			break
		}
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err)
		return
	}
}
