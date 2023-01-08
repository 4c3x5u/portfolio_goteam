package register

import (
	"testing"

	"server/assert"

	"golang.org/x/crypto/bcrypt"
)

// TestPasswordHasher tests both the PasswordHasher and ComparerHash by first hashing
// a inPlaintext string with PasswordHasher and then comparing the original inPlaintext
// to the hash with ComparerHash.
func TestPasswordHasher(t *testing.T) {
	for _, c := range []struct {
		name           string
		inPlaintext    string
		matchPlaintext []byte
		wantErr        error
	}{
		{
			name:           "Match",
			inPlaintext:    "password",
			matchPlaintext: []byte("password"),
			wantErr:        nil,
		},
		{
			name:           "NoMatch",
			inPlaintext:    "password",
			matchPlaintext: []byte("different"),
			wantErr:        bcrypt.ErrMismatchedHashAndPassword,
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
