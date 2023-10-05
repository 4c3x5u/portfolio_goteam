// Package column contains code for responding to HTTP requests made to the
// column API route.
package column

// ReqBody defines the request body for requests handled by Handler.
type ReqBody []struct {
	ID    int `json:"id"`
	Order int `json:"order"`
}

// ResBody defines the response body for the requests handled by Handler.
type ResBody struct {
	Error string `json:"error"`
}
