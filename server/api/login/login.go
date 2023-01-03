// Package login contains code for responding to HTTP requests made to the login
// API route.
package login

// ReqBody defines the request body for Handler.
type ReqBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
