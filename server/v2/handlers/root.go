// Package handlers contains HTTP handlers that are used to response to HTTP
// requests made to server.
package handlers

import (
	"net/http"

	"github.com/kxplxn/goteam/server/v2/relay"
)

// Root is a HTTP handler for the root endpoint.
type Root struct {
	log relay.ErrMsger
}

// NewRoot is the constructor for Root handler.
func NewRoot(log relay.ErrMsger) *Root {
	return &Root{log: log}
}

// ServeHTTP responds to requests made to the root endpoint.
func (h *Root) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	h.log.Msg(w, "app status: OK\navailable endpoints: \n  /register")
}
