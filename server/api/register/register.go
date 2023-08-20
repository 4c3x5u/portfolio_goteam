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
	Msg  string         `json:"message,omitempty"`
	Errs ValidationErrs `json:"errors,omitempty"`
}

// ValidationErrs defines the validation errors returned in ResBody.
type ValidationErrs struct {
	Username []string `json:"username,omitempty"`
	Password []string `json:"password,omitempty"`
}

// Any checks whether there are any validation errors within the ValidationErrs.
func (e ValidationErrs) Any() bool {
	return len(e.Username) > 0 || len(e.Password) > 0
}
