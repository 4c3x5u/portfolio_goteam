package register

// ResBody defines the response body for the register route.
type ResBody struct {
	ErrsValidation *ErrsValidation `json:"errors"`
}

// ErrsValidation defines the structure of error object that can be encoded in
// the register route in the case of an error.
type ErrsValidation struct {
	Username []string `json:"username"`
	Password []string `json:"password"`
}
