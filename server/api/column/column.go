// Package column contains code for responding to HTTP requests made to the
// column API route.
package column

// ResBody defines the response body for the requests handled by Handler.
type ResBody struct {
	Error string `json:"error"`
}
