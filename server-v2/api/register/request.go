package register

// ReqBody defines the request body for the register route.
type ReqBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
