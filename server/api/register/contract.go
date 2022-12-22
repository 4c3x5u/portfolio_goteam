package register

// ReqBody defines the request body for the register route.
type ReqBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ResBody defines the response body for the register route.
type ResBody struct {
	Errs *Errs `json:"errors"`
}

// Errs defines the structure of error object that can be encoded in
// the register route in the case of an error.
type Errs struct {
	Username []string `json:"username"`
	Password []string `json:"password"`
}
