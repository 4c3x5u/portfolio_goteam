package register

import (
	"testing"

	"server/assert"

	"golang.org/x/crypto/bcrypt"
)

// TestPasswordHasher tests the Hash method of the password hasher. It uses
// bcrypt.CompareHashAndPassword to ensure that the result was generated from
// the given plaintext and doesn't match another plaintext string.
func TestPasswordHasher(t *testing.T) {
	for _, c := range []struct {
		name           string
		inPlaintext    string
		matchPlaintext []byte
		wantErr        error
	}{
		{
			name:           "NoMatch",
			inPlaintext:    "password",
			matchPlaintext: []byte("different"),
			wantErr:        bcrypt.ErrMismatchedHashAndPassword,
		},
		{
			name:           "Match",
			inPlaintext:    "password",
			matchPlaintext: []byte("password"),
			wantErr:        nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			sut := NewPasswordHasher()

			hash, err := sut.Hash(c.inPlaintext)

			if err = assert.Nil(err); err != nil {
				t.Error(err)
			}
			err = bcrypt.CompareHashAndPassword(hash, c.matchPlaintext)
			if err = assert.Equal(c.wantErr, err); err != nil {
				t.Error(err)
			}
		})
	}
}
