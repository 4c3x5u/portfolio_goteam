package register

// fakeCreatorUser is a test fake for CreatorUser.
type fakeCreatorUser struct {
	errs *ErrsValidation
}

// CreateUser implements the CreatorUser interface on the fakeCreatorUser
// struct. It returns a pre-set *ErrsValidation object.
func (c *fakeCreatorUser) CreateUser(_, _ string) (*ErrsValidation, error) {
	return c.errs, nil
}

// fakeValidatorReq is a test fake for ValidatorReq
type fakeValidatorReq struct {
	errs *ErrsValidation
}

// Validate implements the ValidatorReq interface on the fakeValidatorReq
// struct.
func (f *fakeValidatorReq) Validate(_ *ReqBody) (_ *ErrsValidation) {
	return f.errs
}

// fakeValidatorStr is a test fake for ValidatorStr.
type fakeValidatorStr struct {
	errs []string
}

// Validate implements the ValidatorStr interface on the fakeValidatorStr
// struct. It returns a pre-set string slice for errsUsername.
func (f *fakeValidatorStr) Validate(_ string) (errs []string) {
	return f.errs
}
