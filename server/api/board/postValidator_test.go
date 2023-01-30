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
		name    string
		reqBody POSTReqBody
		wantErr error
	}{
		{
			name:    "NoName",
			reqBody: POSTReqBody{},
			wantErr: errNameEmpty,
		},
		{
			name:    "NameEmpty",
			reqBody: POSTReqBody{Name: ""},
			wantErr: errNameEmpty,
		},
		{
			name:    "NameTooLong",
			reqBody: POSTReqBody{Name: "boardyboardsyboardkyboardishboardxyz"},
			wantErr: errNameTooLong,
		},
		{
			name:    "IsValid",
			reqBody: POSTReqBody{Name: "My Board"},
			wantErr: nil,
		},
	} {
		err := sut.Validate(c.reqBody)

		if err = assert.Equal(c.wantErr, err); err != nil {
			t.Error(err)
		}
	}
}
