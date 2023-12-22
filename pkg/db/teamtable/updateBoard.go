package teamtable

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/kxplxn/goteam/pkg/db"
)

// BoardUpdater is a type that can be used to update an item in a team's
// boards.
type BoardUpdater struct{ igetput db.DynamoItemGetPutter }

// NewBoardUpdater creates and returns a new BoardUpdater.
func NewBoardUpdater(igetput db.DynamoItemGetPutter) BoardUpdater {
	return BoardUpdater{igetput: igetput}
}

// Update updates a board in the boards of the team with the given ID.
func (d BoardUpdater) Update(
	ctx context.Context, teamID string, board Board,
) error {
	// get the existing team as-is
	out, err := d.igetput.GetItem(ctx, &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: teamID},
		},
		TableName: aws.String(os.Getenv(tableName)),
	})
	if err != nil {
		return err
	}
	if out.Item == nil {
		return db.ErrNoItem
	}

	// unmarshal the team
	var team Team
	if err := attributevalue.UnmarshalMap(out.Item, &team); err != nil {
		return err
	}

	// check that team has boards
	if len(team.Boards) == 0 {
		return db.ErrNoItem
	}

	// check board to be updated exists and update it
	var found bool
	for i, b := range team.Boards {
		if b.ID == board.ID {
			team.Boards[i] = board
			found = true
			break
		}
	}
	if !found {
		return db.ErrNoItem
	}

	// marshal the new team
	newItem, err := attributevalue.MarshalMap(team)
	if err != nil {
		return err
	}

	// update the team based on the new team
	_, err = d.igetput.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      newItem,
		TableName: aws.String(os.Getenv(tableName)),
	})

	return err
}
