package teamtbl

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/kxplxn/goteam/pkg/db"
)

// BoardInserter is a type that can be used to insert an item into a team's
// boards.
type BoardInserter struct{ igetput db.DynamoItemGetPutter }

// NewBoardInserter creates and returns a new BoardInserter.
func NewBoardInserter(igetput db.DynamoItemGetPutter) BoardInserter {
	return BoardInserter{igetput: igetput}
}

// Insert inserts the given board into the boards of the team with the given ID.
func (i BoardInserter) Insert(
	ctx context.Context, teamID string, board Board,
) error {
	// get the existing team as-is
	out, err := i.igetput.GetItem(ctx, &dynamodb.GetItemInput{
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

	// check there are boards to delete
	var dupKey bool
	var count int
	for _, b := range team.Boards {
		if b.ID == board.ID {
			dupKey = true
			break
		}
		count++
	}
	if dupKey {
		return db.ErrDupKey
	}
	if count > 2 {
		return db.ErrLimitReached
	}

	// add the new board into the boards of the team
	team.Boards = append(team.Boards, board)

	// marshal back the team
	newItem, err := attributevalue.MarshalMap(team)
	if err != nil {
		return err
	}

	// update the team
	_, err = i.igetput.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      newItem,
		TableName: aws.String(os.Getenv(tableName)),
	})

	return err
}
