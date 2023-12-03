//go:build utest

package team

import (
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

	t.Run("Err", func(t *testing.T) {
		wantErr := errors.New("failed to get team")
		ig.Out = nil
		ig.Err = wantErr

		_, err := sut.Get("")

		assert.ErrIs(t.Fatal, err, wantErr)
	})

	t.Run("NoItem", func(t *testing.T) {
		wantErr := db.ErrNoItem
		ig.Out = nil
		ig.Err = nil

		_, err := sut.Get("")

		assert.ErrIs(t.Fatal, err, wantErr)
	})

	t.Run("OK", func(t *testing.T) {
		wantTeam := Team{
			ID:      "b8c1dd05-f5de-43ba-bb51-a88051099dba",
			Members: []string{"bob123", "bob124"},
			Boards: []Board{
				{ID: "b8c1dd05-f5de-43ba-bb51-a88051099dba", Name: "Board A"},
				{ID: "630ecf76-383b-42b1-be55-1bc8e3a96e98", Name: "Board B"},
			},
		}
		ig.Out = &dynamodb.GetItemOutput{
			Item: map[string]types.AttributeValue{
				"ID": &types.AttributeValueMemberS{Value: wantTeam.ID},
				"Members": &types.AttributeValueMemberL{
					Value: []types.AttributeValue{
						&types.AttributeValueMemberS{
							Value: wantTeam.Members[0],
						},
						&types.AttributeValueMemberS{
							Value: wantTeam.Members[1],
						},
					},
				},
				"Boards": &types.AttributeValueMemberL{
					Value: []types.AttributeValue{
						&types.AttributeValueMemberM{
							Value: map[string]types.AttributeValue{
								"ID": &types.AttributeValueMemberS{
									Value: wantTeam.Boards[0].ID,
								},
								"Name": &types.AttributeValueMemberS{
									Value: wantTeam.Boards[0].Name,
								},
							},
						},
						&types.AttributeValueMemberM{
							Value: map[string]types.AttributeValue{
								"ID": &types.AttributeValueMemberS{
									Value: wantTeam.Boards[1].ID,
								},
								"Name": &types.AttributeValueMemberS{
									Value: wantTeam.Boards[1].Name,
								},
							},
						},
					},
				},
			},
		}
		ig.Err = nil

		team, err := sut.Get("")

		assert.Nil(t.Fatal, err)
		assert.Equal(t.Error, team.ID, wantTeam.ID)
		assert.AllEqual(t.Error, team.Members, wantTeam.Members)
		assert.Equal(t.Error, team.Boards[0].ID, wantTeam.Boards[0].ID)
		assert.Equal(t.Error, team.Boards[0].Name, wantTeam.Boards[0].Name)
		assert.Equal(t.Error, team.Boards[1].ID, wantTeam.Boards[1].ID)
		assert.Equal(t.Error, team.Boards[1].Name, wantTeam.Boards[1].Name)
	})
}
