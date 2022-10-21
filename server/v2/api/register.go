package api

import (
	"encoding/json"
	"fmt"
	"github.com/kxplxn/goteam/server/v2/relay"
	"net/http"
)

// ReqRegister is the request contract for the register endpoint.
type ReqRegister struct {
	Usn string `json:"username"`
	Pwd string `json:"password"`
	Ref string `json:"referrer"`
}

// HandlerRegister is a HTTP handler for the register endpoint.
type HandlerRegister struct {
	errMsger relay.APIErrMsger
}

// NewHandlerRegister is the constructor for HandlerRegister handler.
func NewHandlerRegister(errMsger relay.APIErrMsger) *HandlerRegister {
	return &HandlerRegister{errMsger: errMsger}
}

// ServeHTTP responds to requests made to the to the register endpoint.
func (h *HandlerRegister) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// accept only POST
	if r.Method != "POST" {
		status := http.StatusMethodNotAllowed
		h.errMsger.Err(w, http.StatusText(status), status)
		return
	}

	// read body into map
	req := &ReqRegister{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		status := http.StatusInternalServerError
		h.errMsger.Err(w, http.StatusText(status), status)
	}

	// rly decoded body
	h.errMsger.Msg(w, fmt.Sprintf(
		"usn: %s\npwd: %s\nref: %s\n",
		req.Usn, req.Pwd, req.Ref,
	))
}
