package board

import (
	"context"
	"database/sql"
	"errors"
)

// Board describes the data needed to insert a board into the database. It
// doesn't represent the final record in the board table - see Record for that.
type Board struct {
	name    string
	adminID string
}

// NewBoard creates and returns a new Board.
func NewBoard(name string, adminID string) Board {
	return Board{name: name, adminID: adminID}
}

// Inserter can be used to create a new record in the board table.
type Inserter struct{ db *sql.DB }

// NewInserter creates and returns a new Inserter.
func NewInserter(db *sql.DB) Inserter { return Inserter{db: db} }

// Insert creates a new record in the board table.
func (i Inserter) Insert(board Board) error {
	// Begin transaction with new empty context.
	ctx := context.Background()
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Insert the new board into the board table.
	var boardID int64
	err = tx.QueryRowContext(
		ctx, "INSERT INTO app.board(name) VALUES ($1) RETURNING id", board.name,
	).Scan(&boardID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Join(err, rollbackErr)
		}
		return err
	}

	// Every time a board is created, the user who creates it must be assigned
	// to it as its admin, and 4 columns must be assigned to the board.

	// Insert a record into the user_board table with the given user and board
	// ID, and an isAdmin field of true (i.e. make the user admin of the board).
	if _, err = tx.ExecContext(
		ctx,
		"INSERT INTO app.user_board(username, boardID, isAdmin) "+
			"VALUES($1, $2, TRUE)",
		board.adminID,
		boardID,
	); err != nil {
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

	return tx.Commit()
}
