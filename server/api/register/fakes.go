package register

// fakeValidatorReq is a test fake for Validator
type fakeValidatorReq struct {
	inReqBody *Req
	outErrs   *Errs
}

// Validate implements the Validator interface on the fakeValidatorReq
// struct.
func (f *fakeValidatorReq) Validate(reqBody *Req) *Errs {
	f.inReqBody = reqBody
	return f.outErrs
}

// fakeValidatorStr is a test fake for ValidatorStr.
type fakeValidatorStr struct {
	inArg   string
	outErrs []string
}

// Validate implements the ValidatorStr interface on the fakeValidatorStr
// struct. It returns a pre-set string slice for errsUsername.
func (f *fakeValidatorStr) Validate(val string) (errs []string) {
	f.inArg = val
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

type fakeHasherPwd struct {
	inPlaintext string
	outHash     []byte
	outErr      error
}

func (f *fakeHasherPwd) Hash(plaintext string) ([]byte, error) {
	f.inPlaintext = plaintext
	return f.outHash, f.outErr
}

type fakeCreatorUser struct {
	inArgs []any
	outErr error
}

func (f *fakeCreatorUser) Create(args ...any) error {
	f.inArgs = args
	return f.outErr
}

type fakeKeeperSession struct {
	inArgs []any
	outErr error
}

func (f *fakeKeeperSession) Create(args ...any) error {
	f.inArgs = args
	return f.outErr
}
