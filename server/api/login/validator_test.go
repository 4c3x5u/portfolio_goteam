//go:build utest

package login

import (
	"testing"

	"github.com/kxplxn/goteam/server/assert"
)

// TestValidator tests the Validate method of Validator to assert that it
// returns the correct boolean value based on the fields of the ReqBody.
func TestValidator(t *testing.T) {
	sut := NewValidator()

	for _, c := range []struct {
		name    string
		reqBody POSTReq
		wantOK  bool
	}{
		{
			name:    "UsernameEmpty",
			reqBody: POSTReq{Username: "", Password: "asdqwe123"},
			wantOK:  false,
		},
		{
			name:    "PasswordEmpty",
			reqBody: POSTReq{Username: "bob123", Password: ""},
			wantOK:  false,
		},
		{
			name:    "UsernameAndPasswordEmpty",
			reqBody: POSTReq{Username: "", Password: ""},
			wantOK:  false,
		},
		{
			name:    "IsValid",
			reqBody: POSTReq{Username: "bob123", Password: "asdqwe123"},
			wantOK:  true,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			ok := sut.Validate(c.reqBody)

			assert.Equal(t.Error, c.wantOK, ok)
		})
	}
}
