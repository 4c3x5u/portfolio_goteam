package register

import "golang.org/x/crypto/bcrypt"

// Hasher describes a type that is used to hash a plaintext string and return
// the hashed bytes.
type Hasher interface{ Hash(string) ([]byte, error) }

// PasswordHasher can be used to hash a given password.
type PasswordHasher struct{}

// NewPasswordHasher creates and returns a new password hasher.
func NewPasswordHasher() PasswordHasher { return PasswordHasher{} }

// Hash hashes a string password and returns the hashed bytes.
func (h PasswordHasher) Hash(plaintext string) ([]byte, error) {
	// https://security.stackexchange.com/questions/17207/recommended-of-rounds-for-bcrypt
	// A cost of 11 causes roughly 20,000 rounds, resulting in about 250ms of
	// compute to generate the hash.
	// TODO: verify the above when live
	return bcrypt.GenerateFromPassword([]byte(plaintext), 11)
}
