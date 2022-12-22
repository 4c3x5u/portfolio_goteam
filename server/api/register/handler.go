// Package register contains types and functions required for the register API
// route (/register).
package register

import (
	"encoding/json"
	"net/http"

	"server/relay"
)

// Handler is a HTTP handler for the register route.
type Handler struct {
	creatorUser CreatorUser
	validator   ValidatorReq
}

// NewHandler is the constructor for Handler.
func NewHandler(creatorUser CreatorUser, validator ValidatorReq) *Handler {
	return &Handler{creatorUser: creatorUser, validator: validator}
}

const errHandlerUsernameTaken = "Username is already taken."

// ServeHTTP responds to requests made to the register route.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	req, res := &ReqBody{}, &ResBody{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		relay.ServerErr(w, err.Error())
		return
	}
	if errs := h.validator.Validate(req); errs != nil {
		res.Errs = errs
		relay.ClientJSON(w, res, http.StatusBadRequest)
		return
	}

	switch err := h.creatorUser.CreateUser(req.Username, req.Password); err {
	case nil:
		return
	case errCreatorUsernameTaken:
		res.Errs = &Errs{Username: []string{errHandlerUsernameTaken}}
		relay.ClientJSON(w, res, http.StatusBadRequest)
		return
	default:
		relay.ServerErr(w, err.Error())
		return
	}
}
