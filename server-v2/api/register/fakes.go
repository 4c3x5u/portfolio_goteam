package register

// fakeValidatorReq is a test fake for ValidatorReq
type fakeValidatorReq struct {
	inReqBody *ReqBody
	outErrs   *Errs
}

// Validate implements the ValidatorReq interface on the fakeValidatorReq
// struct.
func (f *fakeValidatorReq) Validate(reqBody *ReqBody) *Errs {
	f.inReqBody = reqBody
	return f.outErrs
}

// fakeValidatorStr is a test fake for ValidatorStr.
type fakeValidatorStr struct {
	inVal   string
	outErrs []string
}

// Validate implements the ValidatorStr interface on the fakeValidatorStr
// struct. It returns a pre-set string slice for errsUsername.
func (f *fakeValidatorStr) Validate(val string) (errs []string) {
	f.inVal = val
	return f.outErrs
}

// fakeCreatorUser is a test fake for CreatorUser.
type fakeCreatorUser struct {
	inUsername         string
	inPassword         string
	outUsernameIsTaken bool
}

// CreateUser implements the CreatorUser interface on the fakeCreatorUser
// struct. It returns a pre-set *Errs object.
func (f *fakeCreatorUser) CreateUser(username, password string) (bool, error) {
	f.inUsername, f.inPassword = username, password
	return f.outUsernameIsTaken, nil
}
