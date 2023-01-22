package db

import (
	"context"
	"database/sql"
)

// Board represents a record in the board table.
type Board struct {
	name    string
	adminID string
}

// NewBoard creates and returns a new Board.
func NewBoard(name string, adminID string) Board {
	return Board{name: name, adminID: adminID}
}

// BoardInserter can be used to create a new record in the board table.
type BoardInserter struct{ db *sql.DB }

// NewBoardInserter creates and returns a new BoardInserter.
func NewBoardInserter(db *sql.DB) BoardInserter { return BoardInserter{db: db} }

// Insert creates a new record in the board table.
func (i BoardInserter) Insert(board Board) error {
	// Begin transaction with new empty context.
	ctx := context.Background()
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Insert the new board into the board table.
	if _, err = tx.ExecContext(
		ctx, "INSERT INTO app.board(name) VALUES ($1)", board.name,
	); err != nil {
		// TODO: return both errors?
		if errRollback := tx.Rollback(); errRollback != nil {
			return errRollback
		}
		return err
	}

	// Insert a column into the user_board table to set the admin of the
	// newly-created board as the user.
	if _, err = tx.ExecContext(
		ctx,
		"INSERT INTO app.user_board(userID, boardID, isAdmin) VALUES($1, $2, $3)",
		board.adminID,
		board.name,
		true,
	); err != nil {
		// TODO: return both errors?
		if errRollback := tx.Rollback(); errRollback != nil {
			return errRollback
		}
		return err
	}

	return nil
}
