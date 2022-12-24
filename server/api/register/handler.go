// Package register contains types and functions required for the register API
// route (/register).
package register

import (
	"encoding/json"
	"net/http"

	"server/db"

	"server/relay"
)

// Handler is a HTTP handler for the register route.
type Handler struct {
	validator   Validator
	existorUser db.Existor
	hasherPwd   Hasher
	creatorUser db.Creator
}

// NewHandler is the constructor for Handler.
func NewHandler(
	validator Validator,
	existorUser db.Existor,
	hasherPwd Hasher,
	creatorUser db.Creator,
) *Handler {
	return &Handler{
		validator:   validator,
		existorUser: existorUser,
		hasherPwd:   hasherPwd,
		creatorUser: creatorUser,
	}
}

const errFieldUsernameTaken = "Username is already taken."

// ServeHTTP responds to requests made to the register route.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

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

	if userExists, err := h.existorUser.Exists(req.Username); err != nil {
		relay.ServerErr(w, err.Error())
		return
	} else if userExists {
		res.Errs = &Errs{Username: []string{errFieldUsernameTaken}}
		relay.ClientJSON(w, res, http.StatusBadRequest)
		return
	}

	if pwdHash, err := h.hasherPwd.Hash(req.Password); err != nil {
		relay.ServerErr(w, err.Error())
		return
	} else if err := h.creatorUser.Create(req.Username, pwdHash); err != nil {
		relay.ServerErr(w, err.Error())
		return
	}

	relay.ClientJSON(w, res, http.StatusOK)
}
