package register

// ReqBody defines the request body for the register route.
type ReqBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ResBody defines the response body for the register route.
type ResBody struct {
	Msg  string `json:"message"`
	Errs *Errs  `json:"errors"`
}

// Errs defines the structure of request field errors.
type Errs struct {
	Username []string `json:"username"`
	Password []string `json:"password"`
	Session  string   `json:"session"`
}
