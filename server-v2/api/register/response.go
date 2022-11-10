package register

// ResBody defines the response body for the register endpoint.
type ResBody struct {
	Errs *Errs `json:"errors"`
}

// MapErrs implements the ErrsMapper interface from Package test for the test
// cases to be able to inspect the errors object through dynamically set fields.
func (r *ResBody) MapErrs() map[string][]string {
	errsMap := make(map[string][]string)
	errsMap["username"] = r.Errs.Username
	errsMap["password"] = r.Errs.Password
	return errsMap
}

// Errs defines the structure of error object that can be encoded in the
// register endpoint in the case of an error.
type Errs struct {
	Username []string `json:"username"`
	Password []string `json:"password"`
}
