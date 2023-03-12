package midware

import (
	"net/http"
)

// CORS is a middleware that can be used to set up cross-origin resource sharing
// for a given origin.
type CORS struct {
	innerHandler  http.Handler
	allowedOrigin string
}

// NewCORS creates and returns a new CORS.
func NewCORS(innerHandler http.Handler, allowedOrigin string) CORS {
	return CORS{innerHandler: innerHandler, allowedOrigin: allowedOrigin}
}

// ServeHTTP returns access control headers on OPTIONS request. Otherwise,
// it lets the multiplexer handle the request as configured.
func (m CORS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", m.allowedOrigin)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		return
	}
	m.innerHandler.ServeHTTP(w, r)
}
