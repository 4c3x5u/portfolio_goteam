//go:build utest

package board

import (
	"testing"

	"server/assert"
)

// TestPOSTValidato tests the Validate method of POSTValidator to assert that it
// returns the correct boolean value based on the fields of the ReqBody.
func TestPOSTValidator(t *testing.T) {
	sut := NewPOSTValidator()

	for _, c := range []struct {
		name       string
		reqBody    POSTReqBody
		wantErrMsg string
	}{
		{
			name:       "NoName",
			reqBody:    POSTReqBody{},
			wantErrMsg: msgNameEmpty,
		},
		{
			name:       "NameEmpty",
			reqBody:    POSTReqBody{Name: ""},
			wantErrMsg: msgNameEmpty,
		},
		{
			name:       "NameTooLong",
			reqBody:    POSTReqBody{Name: "boardyboardsyboardkyboardishboardxyz"},
			wantErrMsg: msgNameTooLong,
		},
		{
			name:       "IsValid",
			reqBody:    POSTReqBody{Name: "My Board"},
			wantErrMsg: "",
		},
	} {
		msg := sut.Validate(c.reqBody)

		if err := assert.Equal(c.wantErrMsg, msg); err != nil {
			t.Error(msg)
		}
	}
}
