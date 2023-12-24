package loginapi

import (
	"golang.org/x/crypto/bcrypt"
)

// Comparator describes a type that can be used to compare a string with a slice
// of bytes.
type Comparator interface{ Compare([]byte, string) error }

// PasswordComparator is used to compare a plaintext input with a hashed password.
type PasswordComparator struct{}

// NewPasswordComparator creates and returns a new password comparer.
func NewPasswordComparator() PasswordComparator { return PasswordComparator{} }

// Compare compares a plaintext input with a hashed password.
func (c PasswordComparator) Compare(hash []byte, plaintext string) error {
	return bcrypt.CompareHashAndPassword(hash, []byte(plaintext))
}
