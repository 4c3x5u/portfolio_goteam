package dbaccess

import (
	"context"
	"database/sql"
	"errors"
)

// Board represents a record in the board table.
type Board struct {
	id   int
	name string
}

// InBoard describes the data needed to insert a board into the database. It
// doesn't represent the final record in the board table.
type InBoard struct {
	name    string
	adminID string
}

// NewInBoard creates and returns a new InBoard.
func NewInBoard(name string, adminID string) InBoard {
	return InBoard{name: name, adminID: adminID}
}

// BoardSelector can be used to read records from the board table.
type BoardSelector struct{ db *sql.DB }

// NewBoardSelector creates and returns a new BoardSelector.
func NewBoardSelector(db *sql.DB) BoardSelector { return BoardSelector{db: db} }

func (s BoardSelector) Select(id string) (Board, error) {
	var board Board
	err := s.db.
		QueryRow(`SELECT id, name FROM app.board WHERE id = $1`, id).
		Scan(&board.id, &board.name)
	return board, err
}

// BoardInserter can be used to create a new record in the board table.
type BoardInserter struct{ db *sql.DB }

// NewBoardInserter creates and returns a new BoardInserter.
func NewBoardInserter(db *sql.DB) BoardInserter { return BoardInserter{db: db} }

// Insert creates a new record in the board table.
func (i BoardInserter) Insert(board InBoard) error {
	// Begin transaction with new empty context.
	ctx := context.Background()
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Insert the new board into the board table.
	// i.e. Create board.
	var boardID int64
	err = tx.QueryRowContext(
		ctx, "INSERT INTO app.board(name) VALUES ($1) RETURNING id", board.name,
	).Scan(&boardID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return wrapRollbackErr(err, rollbackErr)
		}
		return err
	}

	// Insert a column into the user_board table with the given user and board
	// ID, and an isAdmin field of true (i.e. make the user admin of the board).
	if _, err = tx.ExecContext(
		ctx,
		"INSERT INTO app.user_board(username, boardID, isAdmin) "+
			"VALUES($1, $2, TRUE)",
		board.adminID,
		boardID,
	); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return wrapRollbackErr(err, rollbackErr)
		}
		return err
	}

	tx.Commit()

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
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return wrapRollbackErr(err, rollbackErr)
		}
		return err
	}

	// Delete the row from the board table that corresponds to the board ID.
	// i.e. Delete board.
	if _, err = tx.ExecContext(
		ctx, "DELETE FROM app.board WHERE id = $1", id,
	); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return wrapRollbackErr(err, rollbackErr)
		}
		return err
	}

	tx.Commit()

	return nil
}

// wrapRollbackErr is a helper function to standardise the message for cases
// where a rollback error occurs alongside the error that causes the rollback.
func wrapRollbackErr(err, rollbackErr error) error {
	return errors.New(
		"multiple errors occured:" +
			"\n  (0) err: " + err.Error() +
			"\n  (1) rollbackErr: " + rollbackErr.Error(),
	)
}
