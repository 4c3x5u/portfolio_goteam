package register

// FakeCreatorUser is a test fake for CreatorUser.
type FakeCreatorUser struct {
	errs *ErrsValidation
}

// NewFakeCreatorUser is the constructor for FakeCreatorUser.
func NewFakeCreatorUser(errs *ErrsValidation) *FakeCreatorUser {
	return &FakeCreatorUser{errs: errs}
}

// CreateUser implements the CreatorUser interface on the FakeCreatorUser
// struct. It returns a pre-set *ErrsValidation object.
func (c *FakeCreatorUser) CreateUser(_, _ string) (*ErrsValidation, error) {
	return c.errs, nil
}

// FakeValidatorReq is a test fake for ValidatorReq
type FakeValidatorReq struct {
	errs *ErrsValidation
}

// NewFakeValidatorReq is the constructor for FakeValidatorReq.
func NewFakeValidatorReq() *FakeValidatorReq {
	return &FakeValidatorReq{}
}

// Validate implements the ValidatorReq interface on the FakeValidatorReq
// struct.
func (f *FakeValidatorReq) Validate(_ *ReqBody) (_ *ErrsValidation) {
	return f.errs
}

// FakeValidatorStr is a test fake for ValidatorStr.
type FakeValidatorStr struct {
	errs []string
}

// NewFakeValidatorStr is the constructor for FakeValidatorStr.
func NewFakeValidatorStr() *FakeValidatorStr {
	return &FakeValidatorStr{}
}

// Validate implements the ValidatorStr interface on the FakeValidatorStr
// struct. It returns a pre-set string slice for errsUsername.
func (f *FakeValidatorStr) Validate(_ string) (errs []string) {
	return f.errs
}
