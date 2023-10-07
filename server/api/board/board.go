// Package board contains code for responding to HTTP requests made to the
// board API route.
package board

// ReqBody defines the request body for requests handled by method handlers.
type ReqBody struct {
	Name string `json:"name"`
}

// ResBody defines the response body for requests handled by method handlers.
type ResBody struct {
	Error string `json:"error,omitempty"`
}
