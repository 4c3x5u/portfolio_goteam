// Package handlers contains HTTP handlers that are used to response to HTTP
// requests made to server.
package api

import (
	"net/http"
)

// HandlerRoot is a HTTP handler for the root endpoint.
type HandlerRoot struct {
	log relay.Msger
}

// NewHandlerRoot is the constructor for HandlerRoot handler.
func NewHandlerRoot(msger relay.Msger) *HandlerRoot {
	return &HandlerRoot{log: msger}
}

// ServeHTTP responds to requests made to the root endpoint.
func (h *HandlerRoot) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	h.log.Init(w)
	h.log.Msg("app status: OK\navailable endpoints: \n  /register")
}
