//go:build utest

package loginapi

import (
	"testing"

	"github.com/kxplxn/goteam/pkg/assert"
)

// TestValidator tests the Validate method of Validator to assert that it
// returns the correct boolean value based on the fields of the ReqBody.
func TestValidator(t *testing.T) {
	sut := NewValidator()

	for _, c := range []struct {
		name    string
		reqBody PostReq
		wantOK  bool
	}{
		{
			name:    "UsernameEmpty",
			reqBody: PostReq{Username: "", Password: "asdqwe123"},
			wantOK:  false,
		},
		{
			name:    "PasswordEmpty",
			reqBody: PostReq{Username: "bob123", Password: ""},
			wantOK:  false,
		},
		{
			name:    "UsernameAndPasswordEmpty",
			reqBody: PostReq{Username: "", Password: ""},
			wantOK:  false,
		},
		{
			name:    "IsValid",
			reqBody: PostReq{Username: "bob123", Password: "asdqwe123"},
			wantOK:  true,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			ok := sut.Validate(c.reqBody)

			assert.Equal(t.Error, c.wantOK, ok)
		})
	}
}
