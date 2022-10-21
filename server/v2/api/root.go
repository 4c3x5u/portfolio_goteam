// Package handlers contains HTTP handlers that are used to response to HTTP
// requests made to server.
package api

import (
	"net/http"

	"github.com/kxplxn/goteam/server/v2/relay"
)

// HandlerRoot is a HTTP handler for the root endpoint.
type HandlerRoot struct {
	msger relay.APIMsger
}

// NewHandlerRoot is the constructor for HandlerRoot handler.
func NewHandlerRoot(msger relay.APIMsger) *HandlerRoot {
	return &HandlerRoot{msger: msger}
}

// ServeHTTP responds to requests made to the root endpoint.
func (h *HandlerRoot) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	h.msger.Msg(w, "app status: OK\navailable endpoints: \n  /register")
}
