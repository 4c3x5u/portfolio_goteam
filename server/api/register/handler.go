// Package register contains types and functions required for the register API
// route (/register).
package register

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"time"

	"server/db"

	"server/relay"
)

// Handler is a HTTP handler for the register route.
type Handler struct {
	validator      Validator
	existorUser    db.Existor
	hasherPwd      Hasher
	creatorUser    db.Creator
	creatorSession db.Creator
}

// NewHandler is the constructor for Handler.
func NewHandler(
	validator Validator,
	existorUser db.Existor,
	hasherPwd Hasher,
	creatorUser db.Creator,
	creatorSession db.Creator,
) *Handler {
	return &Handler{
		validator:      validator,
		existorUser:    existorUser,
		hasherPwd:      hasherPwd,
		creatorUser:    creatorUser,
		creatorSession: creatorSession,
	}
}

const errFieldUsernameTaken = "Username is already taken."

// ServeHTTP responds to requests made to the register route.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// only accept post
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// read and validate request
	req, res := &Req{}, &Res{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		relay.ServerErr(w, err.Error())
		return
	}
	if errs := h.validator.Validate(req); errs != nil {
		res.Errs = errs
		relay.ClientJSON(w, res, http.StatusBadRequest)
		return
	}

	// user exists checks/error
	if userExists, err := h.existorUser.Exists(req.Username); err != nil {
		relay.ServerErr(w, err.Error())
		return
	} else if userExists {
		res.Errs = &Errs{Username: []string{errFieldUsernameTaken}}
		relay.ClientJSON(w, res, http.StatusBadRequest)
		return
	}

	// hash password and create user
	if pwdHash, err := h.hasherPwd.Hash(req.Password); err != nil {
		relay.ServerErr(w, err.Error())
		return
	} else if err := h.creatorUser.Create(req.Username, pwdHash); err != nil {
		relay.ServerErr(w, err.Error())
		return
	}

	// keep a new session for this user and set session token cookie
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(1 * time.Hour)
	if err := h.creatorSession.Create(sessionToken, req.Username, expiresAt); err != nil {
		// user successfuly registered but session keeper errored
		res.Errs = &Errs{Session: "register success but session keeper error"}
		relay.ClientErr(w, res, res.Errs.Session, http.StatusUnauthorized)
		return
	} else {
		// register succes, session keeper success, all good...
		http.SetCookie(w, &http.Cookie{
			Name:    "sessionToken",
			Value:   sessionToken,
			Expires: expiresAt,
		})
		relay.ClientJSON(w, res, http.StatusOK)
		return
	}
}
