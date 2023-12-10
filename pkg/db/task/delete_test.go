//go:build utest

package task

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go"

	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/db"
)

func TestDelete(t *testing.T) {
	idel := &db.FakeDynamoItemDeleter{}
	sut := NewDeleter(idel)

	errA := errors.New("failed")
	taskA := Task{
		ID:          "8c5088eb-e86f-4371-86d0-da186dab78a7",
		Title:       "Do something!",
		Description: "Do it!",
		Order:       21,
		Subtasks: []Subtask{
			{Title: "Do a thing", IsDone: true},
			{Title: "Do another thing", IsDone: false},
		},
		BoardID:      "19639b75-45ef-49aa-981e-346c15b0ffbf",
		ColumnNumber: 1,
	}

	for _, c := range []struct {
		name     string
		idelErr  error
		wantTask *Task
		wantErr  error
	}{
		{name: "Err", idelErr: errA, wantTask: nil, wantErr: errA},
		{
			name: "NoItem",
			idelErr: &smithy.OperationError{
				Err: &types.ConditionalCheckFailedException{},
			},
			wantTask: nil,
			wantErr:  db.ErrNoItem,
		},
		{name: "OK", idelErr: nil, wantTask: &taskA, wantErr: nil},
	} {
		t.Run(c.name, func(t *testing.T) {
			idel.Err = c.idelErr

			err := sut.Delete(context.Background(), "")

			assert.Equal(t.Fatal, err, c.wantErr)
		})
	}
}
