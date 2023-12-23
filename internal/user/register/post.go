package register

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/kxplxn/goteam/pkg/cookie"
	"github.com/kxplxn/goteam/pkg/db"
	"github.com/kxplxn/goteam/pkg/db/usertable"
	"github.com/kxplxn/goteam/pkg/log"
)

// PostReq defines the body of POST register requests.
type PostReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// PostResp defines the body of POST register responses.
type PostResp struct {
	Err            string         `json:"error,omitempty"`
	ValidationErrs ValidationErrs `json:"validationErrors,omitempty"`
}

// ValidationErrs defines the validation errors returned in POSTResp.
type ValidationErrs struct {
	Username []string `json:"username,omitempty"`
	Password []string `json:"password,omitempty"`
}

// Any checks whether there are any validation errors within the ValidationErrors.
func (e ValidationErrs) Any() bool {
	return len(e.Username) > 0 || len(e.Password) > 0
}

// PostHandler is a api.MethodHandler that can be used to handle POST register
// requests.
type PostHandler struct {
	reqValidator  ReqValidator
	hasher        Hasher
	inviteDecoder cookie.Decoder[cookie.Invite]
	userInserter  db.Inserter[usertable.User]
	authEncoder   cookie.Encoder[cookie.Auth]
	log           log.Errorer
}

// NewPostHandler creates and returns a new HandlerPost.
func NewPostHandler(
	userValidator ReqValidator,
	inviteDecoder cookie.Decoder[cookie.Invite],
	hasher Hasher,
	userInserter db.Inserter[usertable.User],
	authEncoder cookie.Encoder[cookie.Auth],
	log log.Errorer,
) PostHandler {
	return PostHandler{
		reqValidator:  userValidator,
		hasher:        hasher,
		inviteDecoder: inviteDecoder,
		userInserter:  userInserter,
		authEncoder:   authEncoder,
		log:           log,
	}
}

// ServeHTTP responds to requests made to the register route.
func (h PostHandler) Handle(w http.ResponseWriter, r *http.Request, _ string) {
	// decode request
	var req PostReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// validate request
	vdtErrs := h.reqValidator.Validate(req)
	if vdtErrs.Any() {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(
			PostResp{ValidationErrs: vdtErrs},
		); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	}

	// determine teamID and isAdmin based on invite token.
	ckInvite, err := r.Cookie(cookie.InviteName)
	var teamID string
	var isAdmin bool
	if err == http.ErrNoCookie {
		teamID = uuid.NewString()
		isAdmin = true
	} else if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		invite, err := h.inviteDecoder.Decode(*ckInvite)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(
				PostResp{Err: "Invalid invite token."},
			); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				h.log.Error(err.Error())
			}
			return
		}
		teamID = invite.TeamID
		isAdmin = false
	}

	// hash password
	pwdHash, err := h.hasher.Hash(req.Password)
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// insert a new user into the user table
	if err = h.userInserter.Insert(r.Context(), usertable.NewUser(
		req.Username, pwdHash, isAdmin, teamID,
	)); err == db.ErrDupKey {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(
			PostResp{ValidationErrs: ValidationErrs{
				Username: []string{"Username is already taken."},
			}},
		); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	} else if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// generate an auth token
	ckAuth, err := h.authEncoder.Encode(
		cookie.NewAuth(req.Username, isAdmin, teamID),
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(
			PostResp{
				Err: "You have been registered successfully but something " +
					"went wrong. Please log in using the credentials you " +
					"registered with.",
			},
		); err != nil {
			h.log.Error(err.Error())
		}
		return
	}

	// set auth cookie
	http.SetCookie(w, &ckAuth)
}
