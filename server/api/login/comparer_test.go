package login

import (
	"testing"

	"server/assert"

	"golang.org/x/crypto/bcrypt"
)

func TestComparerHash(t *testing.T) {
	sut := NewComparerHash()

	for _, c := range []struct {
		name        string
		inHashStr   string
		inPlaintext string
		wantIsMatch bool
		wantErr     error
	}{
		{
			name:        "IsMatchTrue",
			inPlaintext: "password",
			inHashStr:   "$2a$04$W4ABZofxx5uoJVgTlYS1wuFHz1LLQaBfoO0iwz/04WWmg9LQdCPsS",
			wantIsMatch: true,
			wantErr:     nil,
		},
		{
			name:        "IsMatchFalse",
			inPlaintext: "password",
			inHashStr:   "$2a$04$ngqMWrzBWyg8KO3MGk1cnOISt3wyeBwbFlkvghSHKBkSYOeO2.7XG",
			wantIsMatch: false,
			wantErr:     nil,
		},
		{
			name:        "ErrTooShort",
			inPlaintext: "password",
			inHashStr:   "$2a$04$ngqMWrzBWyg8K",
			wantIsMatch: false,
			wantErr:     bcrypt.ErrHashTooShort,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			isMatch, err := sut.Compare([]byte(c.inHashStr), c.inPlaintext)
			if err = assert.Equal(c.wantErr, err); err != nil {
				t.Error(err)
			}
			if err = assert.Equal(c.wantIsMatch, isMatch); err != nil {
				t.Error(err)
			}
		})
	}
}
