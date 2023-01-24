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
	// i.e. Create board.
	if _, err = tx.ExecContext(
		ctx, "INSERT INTO app.board(name) VALUES ($1)", board.name,
	); err != nil {
		// TODO: return both errors?
		if errRollback := tx.Rollback(); errRollback != nil {
			return errRollback
		}
		return err
	}

	// Insert a column into the user_board table withe the given user and board
	// ID, and an isAdmin field of true.
	// i.e. Make the user admin of the board.
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

// BoardDeleter can be used to delete a record from the board table as well
// as deleting all of the corresponding relationships records from the
// user_board table.
type BoardDeleter struct{ db *sql.DB }

// NewBoardDeleter creates and returns a new BoardDeleter.
func NewBoardDeleter(db *sql.DB) BoardDeleter { return BoardDeleter{db: db} }

// Delete deletes a record from the board table as well as deleting all of
// the corresponding relationship records from the
// user_board table.
func (d BoardDeleter) Delete(id string) error {
	// Begin transaction with new empty context.
	ctx := context.Background()
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Delete all rows from the user_board table that corresponds to the given
	// board ID.
	// i.e. Delete all of board's relationships.
	if _, err = tx.ExecContext(
		ctx, "DELETE FROM app.user_board WHERE boardID = $1", id,
	); err != nil {
		// TODO: return both errors?
		if errRollback := tx.Rollback(); errRollback != nil {
			return errRollback
		}
		return err
	}

	// Delete the row from the board table that corresponds to the board ID.
	// i.e. Delete board.
	if _, err = tx.ExecContext(
		ctx, "DELETE FROM app.user WHERE boardID = $1", id,
	); err != nil {
		// TODO: return both errors?
		if errRollback := tx.Rollback(); errRollback != nil {
			return errRollback
		}
		return err
	}

	return nil
}
