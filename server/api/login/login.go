// Package login contains types and functions used for serving the login API
// route.
package login

// ReqBody defines the request body for Handler.
type ReqBody struct {
	Username string `json:"username"`
}
