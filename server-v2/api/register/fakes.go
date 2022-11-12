package register

// FakeCreatorUser is a fake object that is meant to be used in tests where
// a CreatorUser is required.
type FakeCreatorUser struct {
	errs *ErrsValidation
}

// NewFakeCreatorUser is the constructor for FakeCreatorUser.
func NewFakeCreatorUser(errs *ErrsValidation) *FakeCreatorUser {
	return &FakeCreatorUser{errs: errs}
}

// CreateUser implements the CreatorUser interface on the FakeCreatorUser type.
// It returns a pre-set *ErrsValidation object.
func (c *FakeCreatorUser) CreateUser(_, _ string) (*ErrsValidation, error) {
	return c.errs, nil
}
