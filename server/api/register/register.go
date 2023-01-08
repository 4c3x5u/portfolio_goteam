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
	Msg            string         `json:"message"`
	ValidationErrs ValidationErrs `json:"errors"`
}

// ValidationErrs defines the validation errors returned in ResBody.
type ValidationErrs struct {
	Username []string `json:"username"`
	Password []string `json:"password"`
	Auth     string   `json:"auth"`
}

// Any checks whether there are any validation errors within the ValidationErrs.
func (e ValidationErrs) Any() bool {
	if len(e.Username) > 0 || len(e.Password) > 0 || e.Auth != "" {
		return true
	}
	return false
}
