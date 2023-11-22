package user

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
)

// InRecord represents a record in the user table.
type InRecord struct {
	TeamID   int
	Username string
	Password []byte
}

// NewInRecord creates and returns a new InRecord.
func NewInRecord(teamID int, username string, password []byte) InRecord {
	return InRecord{TeamID: teamID, Username: username, Password: password}
}

// Inserter can be used to create a new record in the user table.
type Inserter struct{ db *sql.DB }

// NewInserter creates and returns a new Inserter.
func NewInserter(db *sql.DB) Inserter { return Inserter{db: db} }

// Insert creates a new record in the user table.
func (i Inserter) Insert(user InRecord) error {
	ctx := context.Background()
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// If team ID is empty, a new team will be created and the user will be
	// made the admin of that team.
	isAdmin := user.TeamID == -1
	if isAdmin {
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
			`VALUES ($1, $2)`,
		user.Username,
		string(user.Password),
		user.TeamID,
		isAdmin,
	); err != nil {
		return err
	}

	return tx.Commit()
}
