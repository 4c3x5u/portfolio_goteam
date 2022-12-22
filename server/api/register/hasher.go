package register

import (
	"golang.org/x/crypto/bcrypt"
)

type Hasher interface {
	Hash(string) ([]byte, error)
}

type HasherPwd struct{}

func NewHasherPwd() *HasherPwd {
	return &HasherPwd{}
}

func (h *HasherPwd) Hash(plaintext string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.MinCost)
}

type Comparer interface {
	Compare([]byte, string) (bool, error)
}

type ComparerHash struct{}

func (c *ComparerHash) Compare(hash []byte, plaintext string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(hash, []byte(plaintext)); err != nil {
		if err.Error() == "crypto/bcrypt: hashedPassword is not the hash of the given password" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
