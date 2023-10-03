package board

import (
	"server/assert"
	"testing"
)

// TestNameValidator tests the Validate method of NameValidator to assert that
// it returns the correct error message based on board name it's given.
func TestNameValidator(t *testing.T) {
	sut := NewNameValidator()

	msgEmpty := "Board name cannot be empty."

	for _, c := range []struct {
		name       string
		reqBody    ReqBody
		wantErrMsg string
	}{
		{
			name:       "NoName",
			reqBody:    ReqBody{},
			wantErrMsg: msgEmpty,
		},
		{
			name:       "NameEmpty",
			reqBody:    ReqBody{Name: ""},
			wantErrMsg: msgEmpty,
		},
		{
			name: "NameTooLong",
			reqBody: ReqBody{
				Name: "boardyboardsyboardkyboardishboardxyz",
			},
			wantErrMsg: "Board name cannot be longer than 35 characters.",
		},
		{
			name:       "IsValid",
			reqBody:    ReqBody{Name: "My Board"},
			wantErrMsg: "",
		},
	} {
		err := sut.Validate(c.reqBody.Name)

		if c.wantErrMsg == "" {
			if assertErr := assert.Nil(err); assertErr != nil {
				t.Error(assertErr)
			}
		} else {
			if assertErr := assert.Equal(
				c.wantErrMsg, err.Error(),
			); assertErr != nil {
				t.Error(assertErr)
			}
		}
	}
}

// TestIDValidator tests the Validate method of IDValidator to assert
// that it returns the correct error message based on the board ID it's given.
func TestIDValidator(t *testing.T) {
	sut := NewIDValidator()

	for _, c := range []struct {
		name       string
		boardID    string
		wantErrMsg string
	}{
		{
			name:       "Nil",
			boardID:    "",
			wantErrMsg: "Board ID cannot be empty.",
		},
		{
			name:       "NotInt",
			boardID:    "My Board",
			wantErrMsg: "Board ID must be an integer.",
		},
		{
			name:       "Success",
			boardID:    "12",
			wantErrMsg: "",
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			err := sut.Validate(c.boardID)

			if c.wantErrMsg == "" {
				if assertErr := assert.Nil(err); assertErr != nil {
					t.Error(assertErr)
				}
			} else {
				if assertErr := assert.Equal(
					c.wantErrMsg, err.Error(),
				); assertErr != nil {
					t.Error(assertErr)
				}
			}
		})
	}
}
