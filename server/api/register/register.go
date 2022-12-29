// Package register contains types and functions used for serving the register
// API route.
package register

// ReqBody defines the request body for Handler.
type ReqBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ResBody defines the response body for Handler.
type ResBody struct {
	Msg  string `json:"message"`
	Errs *Errs  `json:"errors"`
}

// Errs defines the structure errors returned in ResBody.
type Errs struct {
	Username []string `json:"username"`
	Password []string `json:"password"`
	Session  string   `json:"session"`
}
