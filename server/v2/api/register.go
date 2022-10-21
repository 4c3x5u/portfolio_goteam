package api

import (
	"encoding/json"
	"fmt"
	"github.com/kxplxn/goteam/server/v2/relay"
	"net/http"
)

// ReqRegister is the request contract for the register endpoint.
type ReqRegister struct {
	Usn string `json:"usn"` // username
	Pwd string `json:"pwd"` // password
	Ref string `json:"ref"` // referrer
}

// HandlerRegister is a HTTP handler for the register endpoint.
type HandlerRegister struct {
	log relay.APIErrMsger
}

// NewHandlerRegister is the constructor for HandlerRegister handler.
func NewHandlerRegister(errMsger relay.APIErrMsger) *HandlerRegister {
	return &HandlerRegister{log: errMsger}
}

// ServeHTTP responds to requests made to the to the register endpoint.
func (h *HandlerRegister) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// accept only POST
	if r.Method != "POST" {
		status := http.StatusMethodNotAllowed
		h.log.Err(w, http.StatusText(status), status)
		return
	}

	// decode body into request type
	req := &ReqRegister{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		status := http.StatusInternalServerError
		h.log.Err(w, http.StatusText(status), status)
	}

	// todo: check if user exists in he database

	// relay request fields
	h.log.Msg(w, fmt.Sprintf(
		"usn: %s\npwd: %s\nref: %s\n",
		req.Usn, req.Pwd, req.Ref,
	))
}
