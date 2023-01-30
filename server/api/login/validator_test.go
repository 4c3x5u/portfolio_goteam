package login

import (
	"testing"

	"server/assert"
)

// TestValidator tests the Validate method of Validator to assert that it
// returns the correct boolean value based on the fields of the ReqBody.
func TestValidator(t *testing.T) {
	sut := NewValidator()

	for _, c := range []struct {
		name    string
		reqBody ReqBody
		wantOK  bool
	}{
		{
			name:    "UsernameEmpty",
			reqBody: ReqBody{Username: "", Password: "asdqwe123"},
			wantOK:  false,
		},
		{
			name:    "PasswordEmpty",
			reqBody: ReqBody{Username: "bob123", Password: ""},
			wantOK:  false,
		},
		{
			name:    "UsernameAndPasswordEmpty",
			reqBody: ReqBody{Username: "", Password: ""},
			wantOK:  false,
		},
		{
			name:    "NoEmpty",
			reqBody: ReqBody{Username: "bob123", Password: "asdqwe123"},
			wantOK:  true,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			ok := sut.Validate(c.reqBody)

			if err := assert.Equal(c.wantOK, ok); err != nil {
				t.Error(err)
			}
		})
	}
}
