//go:build utest

package taskapi

import (
	"testing"

	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/db/tasktbl"
	"github.com/kxplxn/goteam/pkg/validator"
)

// TestValidatePostReq tests the ValidatePostReq function to assert that it
// returns the correct error based on the PostReq input.
func TestValidatePostReq(t *testing.T) {
	sut := ValidatePostReq

	for _, c := range []struct {
		name    string
		req     PostReq
		wantErr error
	}{
		{
			name:    "BoardIDEmpty",
			req:     PostReq{BoardID: ""},
			wantErr: errBoardIDEmpty,
		},
		{
			name:    "BoardIDInvalid",
			req:     PostReq{BoardID: "invalid"},
			wantErr: errParseBoardID,
		},
		{
			name: "ColNoTooSmall",
			req: PostReq{
				BoardID: "00000000-0000-0000-0000-000000000000",
				ColNo:   0,
			},
			wantErr: errColNoOutOfBounds,
		},
		{
			name: "ColNoTooBig",
			req: PostReq{
				BoardID: "00000000-0000-0000-0000-000000000000",
				ColNo:   5,
			},
			wantErr: errColNoOutOfBounds,
		},
		{
			name: "TitleEmpty",
			req: PostReq{
				BoardID: "00000000-0000-0000-0000-000000000000",
				ColNo:   2,
				Title:   "",
			},
			wantErr: errTitleEmpty,
		},
		{
			name: "TitleTooLong",
			req: PostReq{
				BoardID: "00000000-0000-0000-0000-000000000000",
				ColNo:   2,
				Title:   "asdqweasdqweasdqweasdqweasdqweasdqweasdqweasdqweasd",
			},
			wantErr: errTitleTooLong,
		},
		{
			name: "DescriptionTooLong",
			req: PostReq{
				BoardID: "00000000-0000-0000-0000-000000000000",
				ColNo:   2,
				Title:   "Some Task",
				Description: "asdqweasdqweasdqweasdqweasdqweasdqweasdqweasdqw" +
					"easdasdqweasdqweasdqweasdqweasdqweasdqweasdqweasdqweasda" +
					"sdqweasdqweasdqweasdqweasdqweasdqweasdqweasdqweasdasdqwe" +
					"asdqweasdqweasdqweasdqweasdqweasdqweasdqweasdasdqweasdqw" +
					"easdqweasdqweasdqweasdqweasdqweasdqweasdasdqweasdqweasdq" +
					"weasdqweasdqweasdqweasdqweasdqweasdasdqweasdqweasdqweasd" +
					"qweasdqweasdqweasdqweasdqweasdasdqweasdqweasdqweasdqweas" +
					"dqweasdqweasdqweasdqweasdasdqweasdqweasdqweasdqweasdqwea" +
					"sdqweasdqweasdqweasdasdqweasdqweasdqweasdqweasdqweasdqwe" +
					"asdqwe",
			},
			wantErr: errDescTooLong,
		},
		{
			name: "SubtaskTitleEmpty",
			req: PostReq{
				BoardID:     "00000000-0000-0000-0000-000000000000",
				ColNo:       2,
				Title:       "Some Task",
				Description: "Some Description",
				Subtasks:    []tasktbl.Subtask{{Title: ""}},
			},
			wantErr: errSubtaskTitleEmpty,
		},
		{
			name: "SubtaskTitleTooLong",
			req: PostReq{
				BoardID:     "00000000-0000-0000-0000-000000000000",
				ColNo:       2,
				Title:       "Some Task",
				Description: "Some Description",
				Subtasks: []tasktbl.Subtask{
					{
						Title: "asdqweasdqweasdqweasdqweasdqweasdqweasdqweasd" +
							"qweasd",
					},
				},
			},
			wantErr: errSubtaskTitleTooLong,
		},
		{
			name: "OrderNegative",
			req: PostReq{
				BoardID:     "00000000-0000-0000-0000-000000000000",
				ColNo:       2,
				Title:       "Some Task",
				Description: "Some Description",
				Subtasks:    []tasktbl.Subtask{{Title: "Some Subtask"}},
				Order:       -1,
			},
			wantErr: errOrderNegative,
		},
		{
			name: "OK",
			req: PostReq{
				BoardID:     "00000000-0000-0000-0000-000000000000",
				ColNo:       2,
				Title:       "Some Task",
				Description: "Some Description",
				Subtasks:    []tasktbl.Subtask{{Title: "Some Subtask"}},
				Order:       0,
			},
			wantErr: nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			err := sut(c.req)
			assert.ErrIs(t.Error, err, c.wantErr)
		})
	}
}

// TestTitleValidator tests the TitleValidator.Validate method.
func TestTitleValidator(t *testing.T) {
	sut := NewTitleValidator()

	for _, c := range []struct {
		name    string
		title   string
		wantErr error
	}{
		{
			name:    "TitleEmpty",
			title:   "",
			wantErr: validator.ErrEmpty,
		},
		{
			name:    "TitleTooLong",
			title:   "asdqweasdqweasdqweasdqweasdqweasdqweasdqweasdqweasd",
			wantErr: validator.ErrTooLong,
		},
		{
			name:    "Success",
			title:   "Some Task",
			wantErr: nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			err := sut.Validate(c.title)
			assert.ErrIs(t.Error, err, c.wantErr)
		})
	}
}
