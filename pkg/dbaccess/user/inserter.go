package user

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
)

// Inserter can be used to create a new record in the user table.
type Inserter struct{ db *sql.DB }

// NewInserter creates and returns a new Inserter.
func NewInserter(db *sql.DB) Inserter { return Inserter{db: db} }

// Insert creates a new record in the user table.
func (i Inserter) Insert(user Record) error {
	ctx := context.Background()
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if user.TeamID == -1 && user.IsAdmin {
		if err = tx.QueryRowContext(
			ctx,
			"INSERT INTO app.team(inviteCode) VALUES ($1) RETURNING id",
			uuid.New().String(),
		).Scan(&user.TeamID); err != nil {
			return err
		}
	}

	if _, err = tx.ExecContext(
		ctx,
		`INSERT INTO app."user"(username, password, teamID, isAdmin) `+
			`VALUES ($1, $2, $3, $4)`,
		user.Username,
		string(user.Password),
		user.TeamID,
		user.IsAdmin,
	); err != nil {
		return err
	}

	return tx.Commit()
}
