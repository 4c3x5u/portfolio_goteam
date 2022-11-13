package register

// ResBody defines the response body for the register route.
type ResBody struct {
	ErrsValidation *ErrsValidation `json:"errors"`
}

// ErrsMap implements the test.ErrsMapper interface to enable test cases to
// inspect the errors object through dynamically defined field names.
func (r *ResBody) ErrsMap() map[string][]string {
	errsMap := make(map[string][]string)
	errsMap["password"] = r.ErrsValidation.Username
	errsMap["password"] = r.ErrsValidation.Password
	return errsMap
}

// ErrsValidation defines the structure of error object that can be encoded in
// the register route in the case of an error.
type ErrsValidation struct {
	Username []string `json:"username"`
	Password []string `json:"password"`
}
