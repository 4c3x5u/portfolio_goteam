package register

// CreatorUser defines the signature for a type that creates a user based on
// a password and a password.
type CreatorUser interface {
	CreateUser(username, password string) (*Errs, error)
}

// CreatorDBUser is a type that is used to create a user in the database
type CreatorDBUser struct {
}

// NewCreatorDBUser is the constructor for CreatorDBUser.
func NewCreatorDBUser() *CreatorDBUser {
	return &CreatorDBUser{}
}

// CreateUser creates a new user in the database.
func (c *CreatorDBUser) CreateUser(_, _ string) (*Errs, error) {
	return nil, nil
}
