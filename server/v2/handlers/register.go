package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/kxplxn/goteam/server/v2/requests"
	"net/http"

	"github.com/kxplxn/goteam/server/v2/relay"
)

// Register is a HTTP handler for the register endpoint.
type Register struct {
	log relay.APIErrMsger
}

// NewRegister is the constructor for Register handler.
func NewRegister(log relay.APIErrMsger) *Register {
	return &Register{log: log}
}

// ServeHTTP responds to requests made to the to the register endpoint.
func (h *Register) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// accept only POST
	if r.Method != "POST" {
		status := http.StatusMethodNotAllowed
		h.log.Err(w, http.StatusText(status), status)
		return
	}

	// read body into map
	req := &requests.Register{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		status := http.StatusInternalServerError
		h.log.Err(w, http.StatusText(status), status)
	}

	// rly decoded body
	h.log.Msg(w, fmt.Sprintf(
		"usn: %s\npwd: %s\nref: %s\n",
		req.Usn, req.Pwd, req.Ref,
	))
}
