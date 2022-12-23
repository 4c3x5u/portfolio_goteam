package register

import (
	"server/assert"
	"testing"
)

// TestHasherComparer tests both the HasherPwd and ComparerHash by first hashing
// a plaintext string with HasherPwd and then comparing the original plaintext
// to the hash with ComparerHash.
func TestHasherComparer(t *testing.T) {
	for _, c := range []struct {
		name        string
		plaintext   string
		hashCompare []byte
		compareText string
		wantIsMatch bool
	}{
		{
			name:        "IsMatch",
			plaintext:   "password",
			compareText: "password",
			wantIsMatch: true,
		},
		{
			name:        "IsNoMatch",
			plaintext:   "password",
			compareText: "notpassword",
			wantIsMatch: false,
		},
	} {
		hasher := &HasherPwd{}
		resHash, err := hasher.Hash(c.plaintext)
		assert.Nil(t, err)

		comparer := &ComparerHash{}
		isMatch, err := comparer.Compare(resHash, c.compareText)
		assert.Nil(t, err)

		assert.Equal(t, c.wantIsMatch, isMatch)
	}
}
