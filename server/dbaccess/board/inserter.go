package board

import (
	"context"
	"database/sql"
	"errors"
)

// InRecord describes the data needed to insert a board into the database.
type InRecord struct {
	name   string
	teamID int
}

// NewInRecord creates and returns a new InRecord.
func NewInRecord(name string, teamID int) InRecord {
	return InRecord{name: name, teamID: teamID}
}

// Inserter can be used to create a new record in the board table.
type Inserter struct{ db *sql.DB }

// NewInserter creates and returns a new Inserter.
func NewInserter(db *sql.DB) Inserter { return Inserter{db: db} }

// Insert creates a new record in the board table.
func (i Inserter) Insert(board InRecord) error {
	// Begin transaction with new empty context.
	ctx := context.Background()
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer tx.Rollback()

	// Insert the new board into the board table.
	var boardID int64
	err = tx.QueryRowContext(
		ctx,
		"INSERT INTO app.board(name, teamID) VALUES ($1, $2) RETURNING id",
		board.name,
		board.teamID,
	).Scan(&boardID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Join(err, rollbackErr)
		}
		return err
	}

	// Insert 4 records into the column table with the returned boardID and
	// order values of 1 to 4.
	for order := 1; order < 5; order++ {
		if _, err = tx.ExecContext(
			ctx,
			`INSERT INTO app."column"(boardID, "order") VALUES ($1, $2)`,
			boardID,
			order,
		); err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return errors.Join(err, rollbackErr)
			}
			return err
		}
	}

	// All went well, commit transaction and return err if occurs.
	return tx.Commit()
}
