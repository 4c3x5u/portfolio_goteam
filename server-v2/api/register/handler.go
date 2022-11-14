// Package register contains types and functions required for the register API
// route (/register).
package register

import (
	"encoding/json"
	"net/http"

	"github.com/kxplxn/goteam/server-v2/relay"
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

// ServeHTTP responds to requests made to the register route.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// accept only POST
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// decode body into request object
	req := &ReqBody{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		relay.ServerErr(w, err.Error())
		return
	}

	// create response object to set values into
	res := &ResBody{}

	// validate the request
	if errs := h.validator.Validate(req); errs != nil {
		res.Errs = errs
		relay.ClientJSON(w, res, http.StatusBadRequest)
		return
	}

	// create user
	errsValidation, err := h.creatorUser.CreateUser(req.Username, req.Password)
	if err != nil {
		relay.ServerErr(w, err.Error())
		return
	}
	if errsValidation != nil {
		res.Errs = errsValidation
		relay.ClientJSON(w, res, http.StatusBadRequest)
		return
	}

	relay.ServerErr(w, "not implemented")
	return
}
