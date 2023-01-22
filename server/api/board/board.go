// Package board contains code for responding to HTTP requests made to the
// board API route.
package board

// ReqBody defines the request body for Handler.
type ResBody struct {
	Error string `json:"error"`
}
