// Package auth contains code that handles authentication-related concerns.
package auth

import "time"

// CookieName is the name of the cookie that the auth token is stored in.
const CookieName = "auth-token"

// Duration is the amount of time that an auth token/session should last.
const Duration = 1 * time.Hour

// WWWAuthenticate is a helper function that returns the name and the value of
// the WWW-Authenticate header that must be set when returning a 401
// Unauthorized response.
func WWWAuthenticate() (string, string) {
	return "WWW-Authenticate", "Bearer"
}
