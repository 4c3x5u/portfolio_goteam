// Package api contains HTTP handlers and possibly other types/functions that
// are used to response to HTTP requests to the server.
package api

import (
	"net/http"

	"github.com/kxplxn/goteam/server/v2/relay"
)

// HandlerRoot is a HTTP handler for the "/" (root) endpoint.
type HandlerRoot struct {
	log relay.ErrMsger
}

// NewHandlerRoot is the constructor for HandlerRoot.
func NewHandlerRoot(log relay.ErrMsger) *HandlerRoot {
	return &HandlerRoot{log: log}
}

// ServeHTTP responds to requests made to the "/" (root) endpoint.
func (h *HandlerRoot) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	h.log.Msg(w, "app status: OK\navailable endpoints: \n  /register")
}
