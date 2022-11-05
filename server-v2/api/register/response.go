package register

// ResBody defines the response body for the register endpoint.
type ResBody struct {
	Errs *Errs `json:"errors"`
}

// Errs defines the structure of error object that can be encoded in the
// register endpoint in the case of an error.
type Errs struct {
	Username []string `json:"username"`
	Password []string `json:"password"`
}
