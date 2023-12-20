package team

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/kxplxn/goteam/pkg/db"
)

// BoardDeleter is a type that can be used to delete an item from a team's
// boards.
type BoardDeleter struct {
	iget db.DynamoItemGetter
	iput db.DynamoItemPutter
}

// NewBoardDeleter creates and returns a new BoardDeleter.
func NewBoardDeleter(
	iget db.DynamoItemGetter, iput db.DynamoItemPutter,
) BoardDeleter {
	return BoardDeleter{iget: iget, iput: iput}
}

// Delete deletes the board with the given ID from the team with the given ID.
func (d BoardDeleter) Delete(
	ctx context.Context, teamID string, boardID string,
) error {
	// get the existing team as-is
	out, err := d.iget.GetItem(ctx, &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: teamID},
		},
		TableName: aws.String(tableName),
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

	// check there are boards to delete
	if len(team.Boards) == 0 {
		return db.ErrNoItem
	}

	// check board to be deleted exists and remove it from team's boards
	var found bool
	newTeam := Team{
		ID:      team.ID,
		Members: team.Members,
		Boards:  make([]Board, len(team.Boards)-1),
	}
	if len(team.Boards) == 1 {
		if team.Boards[0].ID == boardID {
			found = true
		}
	} else {
		for i := 0; i < len(team.Boards)-1; i++ {
			if team.Boards[i].ID == boardID {
				found = true
				i--
				continue
			}
			newTeam.Boards[i] = team.Boards[i]
		}
	}
	if !found {
		return db.ErrNoItem
	}

	// marshal the new team
	newItem, err := attributevalue.MarshalMap(newTeam)
	if err != nil {
		return err
	}

	// update the team based on the new team
	_, err = d.iput.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      newItem,
		TableName: aws.String(tableName),
	})

	return err
}
