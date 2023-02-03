//go:build utest

package auth

import (
	"testing"

	"server/assert"
)

// TestBearerTokenReader tests the Read method of BearerTokenReader.
func TestBearerTokenReader(t *testing.T) {
	sut := NewBearerTokenReader()

	for _, c := range []struct {
		name       string
		authHeader string
		wantToken  string
		wantErr    error
	}{
		{name: "InvalidHeader1", authHeader: "Basic ABCDEFG", wantToken: ""},
		{name: "InvalidHeader2", authHeader: "Bear ABCDEFG", wantToken: ""},
		{name: "EmptyToken", authHeader: "Bearer ", wantToken: ""},
		{name: "Success", authHeader: "Bearer ABCDEFG", wantToken: "ABCDEFG"},
	} {
		t.Run(c.name, func(t *testing.T) {
			token := sut.Read(c.authHeader)

			if err := assert.Equal(c.wantToken, token); err != nil {
				t.Error(err)
			}
		})
	}
}
