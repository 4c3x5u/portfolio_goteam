package register

// Req defines the request body for the register route.
type Req struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Res defines the response body for the register route.
type Res struct {
	Errs *Errs `json:"errors"`
}

// Errs defines the structure of request field errors.
type Errs struct {
	Username []string `json:"username"`
	Password []string `json:"password"`
}
