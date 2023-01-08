package login

import "golang.org/x/crypto/bcrypt"

// Comparer descries a type that compares a string with a slice of bytes.
type Comparer interface {
	Compare([]byte, string) (bool, error)
}

// PasswordComparer is used to compare a plaintext password with a hashed password.
type PasswordComparer struct{}

// NewPasswordComparer creates and returns a new password comparer.
func NewPasswordComparer() PasswordComparer { return PasswordComparer{} }

// Compare compares a plaintext password with a hashed password.
func (c PasswordComparer) Compare(hash []byte, plaintext string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hash, []byte(plaintext))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
