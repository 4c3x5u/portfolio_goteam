// Package board contains code for responding to HTTP requests made to the
// board API route.
package board

// POSTReqBody defines the request body for POST requests handled by Handler.
type POSTReqBody struct {
	Name string `json:"name"`
}

// POSTResBody defines the response body for POST requests handled by Handler.
type POSTResBody struct {
	Error string `json:"error"`
}
