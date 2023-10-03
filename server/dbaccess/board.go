package dbaccess

import (
	"context"
	"database/sql"
	"errors"
)

// InBoard describes the data needed to insert a board into the database. It
// doesn't represent the final record in the board table - see Board for that.
type InBoard struct {
	name    string
	adminID string
}

// NewInBoard creates and returns a new InBoard.
func NewInBoard(name string, adminID string) InBoard {
	return InBoard{name: name, adminID: adminID}
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

// Board represents a record in the board table.
type Board struct {
	id   int
	name string
}

// BoardSelector can be used to read records from the board table.
type BoardSelector struct{ db *sql.DB }

// NewBoardSelector creates and returns a new BoardSelector.
func NewBoardSelector(db *sql.DB) BoardSelector { return BoardSelector{db: db} }

// Select selects a record from the board table with the given id.
func (s BoardSelector) Select(id string) (Board, error) {
	var board Board
	err := s.db.
		QueryRow(`SELECT id, name FROM app.board WHERE id = $1`, id).
		Scan(&board.id, &board.name)
	return board, err
}

// BoardUpdater can be used to update the name field of record in the board
// table.
type BoardUpdater struct{ db *sql.DB }

// NewBoardUpdater is the constructor for BoardUpdater.
func NewBoardUpdater(db *sql.DB) BoardUpdater { return BoardUpdater{db: db} }

// Update updates the name field of a record in the board database with a new
// value.
func (u BoardUpdater) Update(id, newName string) error {
	res, err := u.db.Exec(
		"UPDATE app.board SET name = $1 WHERE id = $2", newName, id,
	)
	if err != nil {
		return err
	}
	if rowsAffected, err := res.RowsAffected(); err != nil {
		return err
	} else if rowsAffected == 0 {
		return errors.New("no rows were affected")
	} else if rowsAffected > 1 {
		return errors.New("more than expected rows were affected")
	}
	return nil
}

// BoardDeleter can be used to delete a record from the board table as well
// as deleting the corresponding relationships records from the user_board
// table.
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

	// Delete all records from the user_board table with the given board ID.
	if _, err = tx.ExecContext(
		ctx, "DELETE FROM app.user_board WHERE boardID = $1", id,
	); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Join(err, rollbackErr)
		}
		return err
	}

	// Delete all records from the column table with the given board ID.
	if _, err = tx.ExecContext(
		ctx, `DELETE FROM app."column" WHERE boardID = $1`, id,
	); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Join(err, rollbackErr)
		}
		return err
	}

	// Delete the record from the board table with the given board ID.
	if _, err = tx.ExecContext(
		ctx, "DELETE FROM app.board WHERE id = $1", id,
	); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Join(err, rollbackErr)
		}
		return err
	}

	return tx.Commit()
}
