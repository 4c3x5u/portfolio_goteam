package register

import "golang.org/x/crypto/bcrypt"

// Hasher represents a type that is used to hash a plaintext string value and
// return hashed bytes alongside any error that occur in the process.
type Hasher interface {
	Hash(string) ([]byte, error)
}

// HasherPwd is a type that is used to hash a given user password.
type HasherPwd struct{}

// NewHasherPwd is the constructor for HasherPwd.
func NewHasherPwd() *HasherPwd {
	return &HasherPwd{}
}

// Hash hashes a given user password and returns the hashed bytes alongside
// an error.
func (h *HasherPwd) Hash(plaintext string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.MinCost)
}

// Comparer represents a type that compares the given bytes to the given
// string value. The first return value should be true if they are a match and
// no errors occur during the comparison process.
type Comparer interface {
	Compare([]byte, string) (bool, error)
}

// ComparerHash is used to compare a given hashed bytes with a plaintext string.
// If the hash was originally created from the plaintext value and no errors
// occured during comparison, the first return value is true.
type ComparerHash struct{}

// NewComparerHash is the constructor for ComparerHash.
func NewComparerHash() *ComparerHash {
	return &ComparerHash{}
}

// Compare compares the given hashed bytes with the given plaintext string. The
// first return value communicates whether it was a match. The second return
// value is for any errors that may ocur during comparison.
func (c *ComparerHash) Compare(hash []byte, plaintext string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(hash, []byte(plaintext)); err != nil {
		if err.Error() == "crypto/bcrypt: hashedPassword is not the hash of the given password" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
