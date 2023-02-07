// Package api contains code for serving the various API routes of the app.
// The code is divided into sub-packages that each correspond to a single API
// route.
package api

import "net/http"

// MethodHandler describes a type that can be used to serve a certain part of an
// API route that corresponds to a specific HTTP method. It is intended for the
// Handle method to be called after authentication, and with the authenticated
// user's username as the third argument.
type MethodHandler interface {
	Handle(w http.ResponseWriter, r *http.Request, username string)
}

// AllowedMethods takes in a slice of allowed HTTP methods and returns the key
// and the value for the Access-Control-Allow-Methods header.
func AllowedMethods(methods ...string) (string, string) {
	if len(methods) == 0 {
		return "", ""
	}

	allowedMethods := ""
	for _, m := range methods {
		allowedMethods += m + ", "
	}

	return "Access-Control-Allow-Methods",
		allowedMethods[:len(allowedMethods)-2]
}
