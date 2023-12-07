//go:build utest

package team

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/db"
)

func TestGetter(t *testing.T) {
	ig := &db.FakeItemGetter{}
	sut := NewGetter(ig)

	errA := errors.New("failed to get team")
	teamA := Team{
		ID:      "b8c1dd05-f5de-43ba-bb51-a88051099dba",
		Members: []string{"bob123", "bob124"},
		Boards: []Board{
			{ID: "b8c1dd05-f5de-43ba-bb51-a88051099dba", Name: "Board A"},
			{ID: "630ecf76-383b-42b1-be55-1bc8e3a96e98", Name: "Board B"},
		},
	}

	for _, c := range []struct {
		name     string
		igOut    *dynamodb.GetItemOutput
		igErr    error
		wantTeam *Team
		wantErr  error
	}{
		{
			name:     "Err",
			igOut:    nil,
			igErr:    errA,
			wantTeam: nil,
			wantErr:  errA,
		},
		{
			name:     "NoItem",
			igOut:    nil,
			igErr:    nil,
			wantTeam: nil,
			wantErr:  db.ErrNoItem,
		},
		{
			name: "OK",
			igOut: &dynamodb.GetItemOutput{
				Item: map[string]types.AttributeValue{
					"ID": &types.AttributeValueMemberS{Value: teamA.ID},
					"Members": &types.AttributeValueMemberL{
						Value: []types.AttributeValue{
							&types.AttributeValueMemberS{
								Value: teamA.Members[0],
							},
							&types.AttributeValueMemberS{
								Value: teamA.Members[1],
							},
						},
					},
					"Boards": &types.AttributeValueMemberL{
						Value: []types.AttributeValue{
							&types.AttributeValueMemberM{
								Value: map[string]types.AttributeValue{
									"ID": &types.AttributeValueMemberS{
										Value: teamA.Boards[0].ID,
									},
									"Name": &types.AttributeValueMemberS{
										Value: teamA.Boards[0].Name,
									},
								},
							},
							&types.AttributeValueMemberM{
								Value: map[string]types.AttributeValue{
									"ID": &types.AttributeValueMemberS{
										Value: teamA.Boards[1].ID,
									},
									"Name": &types.AttributeValueMemberS{
										Value: teamA.Boards[1].Name,
									},
								},
							},
						},
					},
				},
			},
			wantTeam: &teamA,
			wantErr:  nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			ig.Out = c.igOut
			ig.Err = c.igErr

			team, err := sut.Get(context.Background(), "")

			assert.ErrIs(t.Fatal, err, c.wantErr)

			if c.wantTeam != nil {
				t.Log(c.wantTeam.ID)

				assert.Equal(t.Error, team.ID, c.wantTeam.ID)
				assert.AllEqual(t.Error, team.Members, c.wantTeam.Members)
				for i, wb := range c.wantTeam.Boards {
					assert.Equal(t.Error, team.Boards[i].ID, wb.ID)
					assert.Equal(t.Error, team.Boards[i].Name, wb.Name)
				}
			}
		})
	}
}
