package register

import (
	"testing"

	"server/assert"

	"golang.org/x/crypto/bcrypt"
)

// TestHasherComparer tests both the HasherPwd and ComparerHash by first hashing
// a inPlaintext string with HasherPwd and then comparing the original inPlaintext
// to the hash with ComparerHash.
func TestHasherComparer(t *testing.T) {
	for _, c := range []struct {
		name           string
		inPlaintext    string
		matchPlaintext string
		wantErr        error
	}{
		{
			name:           "IsMatch",
			inPlaintext:    "password",
			matchPlaintext: "password",
			wantErr:        nil,
		},
		{
			name:           "IsNoMatch",
			inPlaintext:    "password",
			matchPlaintext: "notthesame",
			wantErr:        bcrypt.ErrMismatchedHashAndPassword,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			hasher := &HasherPwd{}
			resHash, err := hasher.Hash(c.inPlaintext)
			assert.Nil(t, err)
			err = bcrypt.CompareHashAndPassword(resHash, []byte(c.matchPlaintext))
			if c.wantErr != nil {
				assert.Equal(t, c.wantErr.Error(), err.Error())
			} else {
				assert.Nil(t, c.wantErr)
			}
		})
	}
}
