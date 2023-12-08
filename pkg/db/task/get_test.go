//go:build utest

package task

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/db"
)

func TestGetter(t *testing.T) {
	ig := &db.FakeItemGetter{}
	sut := NewGetter(ig)

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
		igOut    *dynamodb.GetItemOutput
		igErr    error
		wantTask *Task
		wantErr  error
	}{
		{
			name:     "Err",
			igOut:    nil,
			igErr:    errA,
			wantTask: nil,
			wantErr:  errA,
		},
		{
			name:     "NoItem",
			igOut:    &dynamodb.GetItemOutput{Item: nil},
			igErr:    nil,
			wantTask: nil,
			wantErr:  db.ErrNoItem,
		},
		{
			name: "Err",
			igOut: &dynamodb.GetItemOutput{
				Item: map[string]types.AttributeValue{
					"ID":    &types.AttributeValueMemberS{Value: taskA.ID},
					"Title": &types.AttributeValueMemberS{Value: taskA.Title},
					"Description": &types.AttributeValueMemberS{
						Value: taskA.Description,
					},
					"Order": &types.AttributeValueMemberN{
						Value: strconv.Itoa(taskA.Order),
					},
					"Subtasks": &types.AttributeValueMemberL{
						Value: []types.AttributeValue{
							&types.AttributeValueMemberM{
								Value: map[string]types.AttributeValue{
									"Title": &types.AttributeValueMemberS{
										Value: taskA.Subtasks[0].Title,
									},
									"IsDone": &types.AttributeValueMemberBOOL{
										Value: taskA.Subtasks[0].IsDone,
									},
								},
							},
							&types.AttributeValueMemberM{
								Value: map[string]types.AttributeValue{
									"Title": &types.AttributeValueMemberS{
										Value: taskA.Subtasks[1].Title,
									},
									"IsDone": &types.AttributeValueMemberBOOL{
										Value: taskA.Subtasks[1].IsDone,
									},
								},
							},
						},
					},
					"BoardID": &types.AttributeValueMemberS{
						Value: taskA.BoardID,
					},
					"ColumnNumber": &types.AttributeValueMemberN{
						Value: strconv.Itoa(taskA.ColumnNumber),
					},
				},
			},
			igErr:    nil,
			wantTask: &taskA,
			wantErr:  nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			ig.Out = c.igOut
			ig.Err = c.igErr

			task, err := sut.Get(context.Background(), "")

			assert.Equal(t.Fatal, err, c.wantErr)
			if c.wantTask != nil {
				assert.Equal(t.Error, task.ID, c.wantTask.ID)
				assert.Equal(t.Error, task.Title, c.wantTask.Title)
				assert.Equal(t.Error, task.Description, c.wantTask.Description)
				assert.Equal(t.Error, task.Order, c.wantTask.Order)
				assert.Equal(t.Error, task.BoardID, c.wantTask.BoardID)
				assert.Equal(t.Error,
					task.ColumnNumber, c.wantTask.ColumnNumber,
				)

				for i, wst := range c.wantTask.Subtasks {
					assert.Equal(t.Error, task.Subtasks[i].Title, wst.Title)
					assert.Equal(t.Error, task.Subtasks[i].IsDone, wst.IsDone)
				}
			}
		})
	}
}
