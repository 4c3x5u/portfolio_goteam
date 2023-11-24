// Package column contains code for responding to HTTP requests made to the
// column API route.
package column

// Task represents an item in the body of a request made to the column endpoint.
type Task struct {
	ID    int `json:"id"`
	Order int `json:"order"`
}

// ReqBody defines the request body for requests handled by Handler.
type ReqBody []Task

// ResBody defines the response body for the requests handled by Handler.
type ResBody struct {
	Error string `json:"error"`
}
