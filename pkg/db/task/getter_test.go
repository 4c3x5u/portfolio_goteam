//go:build utest

package task

import (
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

	t.Run("Err", func(t *testing.T) {
		wantErr := errors.New("failed to get item")
		ig.Out = nil
		ig.Err = wantErr

		_, err := sut.Get("")

		assert.Equal(t.Fatal, err, wantErr)
	})

	t.Run("NoItem", func(t *testing.T) {
		wantErr := db.ErrNoItem
		ig.Out = nil
		ig.Err = nil

		_, err := sut.Get("")

		assert.Equal(t.Fatal, err, wantErr)
	})

	t.Run("NoItem", func(t *testing.T) {
		wantTask := Task{
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
		ig.Out = &dynamodb.GetItemOutput{
			Item: map[string]types.AttributeValue{
				"ID": &types.AttributeValueMemberS{Value: wantTask.ID},
				"Title": &types.AttributeValueMemberS{
					Value: wantTask.Title,
				},
				"Description": &types.AttributeValueMemberS{
					Value: wantTask.Description,
				},
				"Order": &types.AttributeValueMemberN{
					Value: strconv.Itoa(wantTask.Order),
				},
				"Subtasks": &types.AttributeValueMemberL{
					Value: []types.AttributeValue{
						&types.AttributeValueMemberM{
							Value: map[string]types.AttributeValue{
								"Title": &types.AttributeValueMemberS{
									Value: wantTask.Subtasks[0].Title,
								},
								"IsDone": &types.AttributeValueMemberBOOL{
									Value: wantTask.Subtasks[0].IsDone,
								},
							},
						},
						&types.AttributeValueMemberM{
							Value: map[string]types.AttributeValue{
								"Title": &types.AttributeValueMemberS{
									Value: wantTask.Subtasks[1].Title,
								},
								"IsDone": &types.AttributeValueMemberBOOL{
									Value: wantTask.Subtasks[1].IsDone,
								},
							},
						},
					},
				},
				"BoardID": &types.AttributeValueMemberS{
					Value: wantTask.BoardID,
				},
				"ColumnNumber": &types.AttributeValueMemberN{
					Value: strconv.Itoa(wantTask.ColumnNumber),
				},
			},
		}
		ig.Err = nil

		task, err := sut.Get("")

		assert.Nil(t.Fatal, err)
		assert.Equal(t.Error, task.ID, wantTask.ID)
		assert.Equal(t.Error, task.Title, wantTask.Title)
		assert.Equal(t.Error, task.Description, wantTask.Description)
		assert.Equal(t.Error, task.Order, wantTask.Order)
		assert.Equal(t.Error,
			task.Subtasks[0].Title, wantTask.Subtasks[0].Title,
		)
		assert.Equal(t.Error,
			task.Subtasks[0].IsDone, wantTask.Subtasks[0].IsDone,
		)
		assert.Equal(t.Error,
			task.Subtasks[1].Title, wantTask.Subtasks[1].Title,
		)
		assert.Equal(t.Error,
			task.Subtasks[1].IsDone, wantTask.Subtasks[1].IsDone,
		)
		assert.Equal(t.Error, task.BoardID, wantTask.BoardID)
		assert.Equal(t.Error, task.ColumnNumber, wantTask.ColumnNumber)
	})
}
