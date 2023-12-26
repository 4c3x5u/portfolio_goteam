//go:build utest

package tasktbl

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

func TestRetrieverByTeam(t *testing.T) {
	queryer := &db.FakeDynamoQueryer{}
	sut := NewRetrieverByTeam(queryer)

	errA := errors.New("failed")
	someTasks := []Task{
		{
			TeamID:      "577965d9-c7ba-4a18-ae7b-47d879b12879",
			ID:          "8c5088eb-e86f-4371-86d0-da186dab78a7",
			BoardID:     "19639b75-45ef-49aa-981e-346c15b0ffbf",
			ColNo:       1,
			Title:       "Do something!",
			Description: "Do it!",
			Order:       21,
			Subtasks: []Subtask{
				{Title: "Do a thing", IsDone: true},
				{Title: "Do another thing", IsDone: false},
			},
		},
		{
			TeamID:      "577965d9-c7ba-4a18-ae7b-47d879b12879",
			ID:          "0c328813-e1b1-4371-86d2-da184d567877",
			BoardID:     "19639b75-45ef-49aa-981e-346c15b0ffbf",
			ColNo:       0,
			Title:       "Do something again!",
			Description: "Dooooooo it!",
			Order:       52,
			Subtasks: []Subtask{
				{Title: "Do a thing again", IsDone: true},
				{Title: "Do another thing again", IsDone: false},
			},
		},
	}

	for _, c := range []struct {
		name      string
		dqOut     *dynamodb.QueryOutput
		dqErr     error
		wantTasks []Task
		wantErr   error
	}{
		{
			name:      "Err",
			dqOut:     nil,
			dqErr:     errA,
			wantTasks: []Task{},
			wantErr:   errA,
		},
		{
			name: "OK",
			dqOut: func() *dynamodb.QueryOutput {
				var out dynamodb.QueryOutput
				for _, t := range someTasks {
					out.Items = append(out.Items, map[string]types.
						AttributeValue{
						"TeamID": &types.AttributeValueMemberN{
							Value: t.TeamID,
						},
						"BoardID": &types.AttributeValueMemberS{
							Value: t.BoardID,
						},
						"ColNo": &types.AttributeValueMemberN{
							Value: strconv.Itoa(t.ColNo),
						},
						"ID":    &types.AttributeValueMemberS{Value: t.ID},
						"Title": &types.AttributeValueMemberS{Value: t.Title},
						"Description": &types.AttributeValueMemberS{
							Value: t.Description,
						},
						"Order": &types.AttributeValueMemberN{
							Value: strconv.Itoa(t.Order),
						},
						"Subtasks": &types.AttributeValueMemberL{
							Value: []types.AttributeValue{
								&types.AttributeValueMemberM{
									Value: map[string]types.AttributeValue{
										"Title": &types.AttributeValueMemberS{
											Value: t.Subtasks[0].Title,
										},
										"IsDone": &types.
											AttributeValueMemberBOOL{
											Value: t.Subtasks[0].IsDone,
										},
									},
								},
								&types.AttributeValueMemberM{
									Value: map[string]types.AttributeValue{
										"Title": &types.
											AttributeValueMemberS{
											Value: t.Subtasks[1].Title,
										},
										"IsDone": &types.
											AttributeValueMemberBOOL{
											Value: t.Subtasks[1].IsDone,
										},
									},
								},
							},
						},
					})
				}
				return &out
			}(),
			dqErr:     nil,
			wantTasks: someTasks,
			wantErr:   nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			queryer.Out = c.dqOut
			queryer.Err = c.dqErr

			tasks, err := sut.Retrieve(context.Background(), "")

			assert.Equal(t.Fatal, err, c.wantErr)
			assert.Equal(t.Error, len(tasks), len(c.wantTasks))
			for i, wt := range c.wantTasks {
				task := tasks[i]

				assert.Equal(t.Error, task.TeamID, wt.TeamID)
				assert.Equal(t.Error, task.ID, wt.ID)
				assert.Equal(t.Error, task.Title, wt.Title)
				assert.Equal(t.Error, task.Description, wt.Description)
				assert.Equal(t.Error, task.Order, wt.Order)
				assert.Equal(t.Error, task.BoardID, wt.BoardID)
				assert.Equal(t.Error,
					task.ColNo, wt.ColNo,
				)

				for j, wst := range wt.Subtasks {
					assert.Equal(t.Error, task.Subtasks[j].Title, wst.Title)
					assert.Equal(t.Error, task.Subtasks[j].IsDone, wst.IsDone)
				}
			}
		})
	}
}
