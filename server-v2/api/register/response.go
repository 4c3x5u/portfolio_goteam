package register

// Res defines the resposne type for the register endpoint.
type Res struct {
	Errs *Errs `json:"errors"`
}

// Errs defines the structure of error object that can be encoded in the
// register endpoint in the case of an error.
type Errs struct {
	Username []string `json:"username"`
}
