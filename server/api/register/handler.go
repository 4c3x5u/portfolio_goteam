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
	existorUser db.Existor
	validator   ValidatorReq
}

// NewHandler is the constructor for Handler.
func NewHandler(existorUser db.Existor, validator ValidatorReq) *Handler {
	return &Handler{existorUser: existorUser, validator: validator}
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

	userExists, err := h.existorUser.Exists(req.Username)
	if err != nil {
		relay.ServerErr(w, err.Error())
		return
	}
	if userExists {
		res.Errs = &Errs{Username: []string{errHandlerUsernameTaken}}
		relay.ClientJSON(w, res, http.StatusBadRequest)
		return
	}

	// todo: create the user
}
