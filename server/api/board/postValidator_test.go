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

	msgEmpty := "Board name cannot be empty."

	for _, c := range []struct {
		name       string
		reqBody    POSTReqBody
		wantErrMsg string
	}{
		{
			name:       "NoName",
			reqBody:    POSTReqBody{},
			wantErrMsg: msgEmpty,
		},
		{
			name:       "NameEmpty",
			reqBody:    POSTReqBody{Name: ""},
			wantErrMsg: msgEmpty,
		},
		{
			name: "NameTooLong",
			reqBody: POSTReqBody{
				Name: "boardyboardsyboardkyboardishboardxyz",
			},
			wantErrMsg: "Board name cannot be longer than 35 characters.",
		},
		{
			name:       "IsValid",
			reqBody:    POSTReqBody{Name: "My Board"},
			wantErrMsg: "",
		},
	} {
		errMsg := sut.Validate(c.reqBody)

		if err := assert.Equal(c.wantErrMsg, errMsg); err != nil {
			t.Error(err)
		}
	}
}
