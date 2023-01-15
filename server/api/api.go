// Package api contains code for serving the various API endpoints of the app.
// The code is divided into sub-packages that each correspond to a single API
// endpoint.
package api

// WWWAuthenticate is a helper function that returns the name and the value of
// the WWW-Authenticate header that must be set when returning a 401
// Unauthorized response.
func WWWAuthenticate() (string, string) {
	return "WWW-Authenticate", "Bearer"
}
