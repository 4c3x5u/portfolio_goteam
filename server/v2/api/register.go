package api

import (
	"encoding/json"
	"fmt"
	"github.com/kxplxn/goteam/server/v2/db"
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
		h.log.ErrCode(w, http.StatusMethodNotAllowed)
		return
	}

	// decode body into request type
	req := &ReqRegister{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.ErrCode(w, http.StatusInternalServerError)
	}

	// todo: check if user exists in he database
	client, ctx, cancel, err := db.Connect("")
	if err != nil {
		h.log.ErrCode(w, http.StatusInternalServerError)
	}

	defer db.Close(client, ctx, cancel)

	// relay request fields
	h.log.Msg(w, fmt.Sprintf(
		"usn: %s\npwd: %s\nref: %s\n",
		req.Usn, req.Pwd, req.Ref,
	))
}
