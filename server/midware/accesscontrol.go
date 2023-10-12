package midware

import (
	"net/http"
)

// AccessControl is a middleware that can be used to set up the access control
// headers needed for a http.Handler.
type AccessControl struct {
	innerHandler  http.Handler
	allowedOrigin string
}

// NewAccessControl creates and returns a new AccessControl.
func NewAccessControl(
	innerHandler http.Handler, allowedOrigin string,
) AccessControl {
	return AccessControl{
		innerHandler:  innerHandler,
		allowedOrigin: allowedOrigin,
	}
}

// ServeHTTP returns access control headers on OPTIONS request. Otherwise,
// it lets the multiplexer handle the request as configured.
func (m AccessControl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", m.allowedOrigin)
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Add("Access-Control-Allow-Credentials", "true")

	if r.Method == http.MethodOptions {
		return
	}

	m.innerHandler.ServeHTTP(w, r)
}
