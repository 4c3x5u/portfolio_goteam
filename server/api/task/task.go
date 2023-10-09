// Package task contains code for responding to HTTP requests made to the task
// API route.
package task

// ResBody defines the response body for requests handled by method handlers.
type ResBody struct {
	Error string `json:"error"`
}
