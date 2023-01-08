package login

import (
	"testing"

	"server/assert"

	"golang.org/x/crypto/bcrypt"
)

func TestPasswordComparer(t *testing.T) {
	sut := NewPasswordComparer()

	for _, c := range []struct {
		name        string
		inHash      []byte
		inPlaintext string
		wantMatch   bool
		wantErr     error
	}{
		{
			name:        "Match",
			inPlaintext: "password",
			inHash:      []byte("$2a$04$W4ABZofxx5uoJVgTlYS1wuFHz1LLQaBfoO0iwz/04WWmg9LQdCPsS"),
			wantMatch:   true,
			wantErr:     nil,
		},
		{
			name:        "NoMatch",
			inPlaintext: "password",
			inHash:      []byte("$2a$04$ngqMWrzBWyg8KO3MGk1cnOISt3wyeBwbFlkvghSHKBkSYOeO2.7XG"),
			wantMatch:   false,
			wantErr:     nil,
		},
		{
			name:        "Error",
			inPlaintext: "password",
			inHash:      []byte("$2a$04$ngqMWrzBWyg8K"),
			wantMatch:   false,
			wantErr:     bcrypt.ErrHashTooShort,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			match, err := sut.Compare(c.inHash, c.inPlaintext)
			if err = assert.Equal(c.wantErr, err); err != nil {
				t.Error(err)
			}
			if err = assert.Equal(c.wantMatch, match); err != nil {
				t.Error(err)
			}
		})
	}
}
