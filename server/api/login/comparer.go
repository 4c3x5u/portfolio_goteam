package login

import "golang.org/x/crypto/bcrypt"

// Comparer describes a type that can be used to compare a string with a slice
// of bytes.
type Comparer interface{ Compare([]byte, string) error }

// PasswordComparer is used to compare a plaintext input with a hashed password.
type PasswordComparer struct{}

// NewPasswordComparer creates and returns a new password comparer.
func NewPasswordComparer() PasswordComparer { return PasswordComparer{} }

// Compare compares a plaintext input with a hashed password.
func (c PasswordComparer) Compare(hash []byte, plaintext string) error {
	return bcrypt.CompareHashAndPassword(hash, []byte(plaintext))
}
