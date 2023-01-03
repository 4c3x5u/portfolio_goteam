package register

import "golang.org/x/crypto/bcrypt"

// Hasher represents a type that is used to hash a plaintext string value and
// return hashed bytes alongside any error that occur in the process.
type Hasher interface{ Hash(string) ([]byte, error) }

// HasherPwd is a type that is used to hash a given user password.
type HasherPwd struct{}

// NewHasherPwd is the constructor for HasherPwd.
func NewHasherPwd() *HasherPwd { return &HasherPwd{} }

// Hash hashes a given user password and returns the hashed bytes alongside
// an error.
func (h *HasherPwd) Hash(plaintext string) ([]byte, error) {
	// https://security.stackexchange.com/questions/17207/recommended-of-rounds-for-bcrypt
	// A cost of 11 causes roughly 20,000 rounds, resulting in about 250ms of
	// compute to generate the hash.
	// TODO: verify the above when live
	return bcrypt.GenerateFromPassword([]byte(plaintext), 11)
}
