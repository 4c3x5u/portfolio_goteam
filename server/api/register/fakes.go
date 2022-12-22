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

// fakeExistorUser is a test fake for Existor.
type fakeExistorUser struct {
	inUsername string
	outExists  bool
	outErr     error
}

// Exists implements the Existor interface on the fakeExistorUser
// struct. It returns a pre-set *Errs object.
func (f *fakeExistorUser) Exists(username string) (bool, error) {
	f.inUsername = username
	return f.outExists, f.outErr
}
