package login

import "golang.org/x/crypto/bcrypt"

// HashComparer represents a type that compares the given bytes to the given
// string value. The first return value should be true if they are a match and
// no errors occur during the comparison process.
type HashComparer interface {
	Compare([]byte, string) (bool, error)
}

// PasswordComparer is used to compare a given plaintext password with a hashed
// password.
type PasswordComparer struct{}

// NewPasswordComparer is the constructor for PasswordComparer.
func NewPasswordComparer() *PasswordComparer { return &PasswordComparer{} }

// Compare compares the given hashed bytes with the given plaintext string. The
// first return value communicates whether it was a match. The second return
// value is for any errors that may ocur during comparison.
func (c *PasswordComparer) Compare(hash []byte, plaintext string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hash, []byte(plaintext))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
