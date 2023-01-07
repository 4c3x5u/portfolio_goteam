// Package register contains code for responding to HTTP requests made to the
// register API route.
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
	Auth     string   `json:"auth"`
}
