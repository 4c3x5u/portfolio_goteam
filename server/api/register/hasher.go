package register

import "golang.org/x/crypto/bcrypt"

// Hasher describes a type that is used to hash a plaintext string and return
// the hashed bytes.
type Hasher interface{ Hash(string) ([]byte, error) }

// PasswordHasher can be used to hash a given password.
type PasswordHasher struct{}

// NewPasswordHasher creates and returns a new PasswordHasher.
func NewPasswordHasher() PasswordHasher { return PasswordHasher{} }

// Hash hashes a string password and returns the hashed bytes.
func (h PasswordHasher) Hash(plaintext string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(plaintext), 11)
}
