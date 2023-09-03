// Package board contains code for responding to HTTP requests made to the
// board API route.
package board

// ReqBody defines the request body for requests handled by Handler.
type ReqBody struct {
	Name string `json:"name"`
}

// ResBody defines the response body for requests handled by Handler.
type ResBody struct {
	Error string `json:"error,omitempty"`
}
