package token

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// InviteName is the name of the invite token.
const InviteName = "invite-token"

// Invite defines the body of an Invite token.
type Invite struct{ TeamID string }

// NewInvite creates and returns a new Invite.
func NewInvite(teamID string) Invite {
	return Invite{TeamID: teamID}
}

// Encode encodes an Invite into a JWT string.
func EncodeInvite(exp time.Time, inv Invite) (string, error) {
	tk, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"teamID": inv.TeamID,
		"exp":    exp.Unix(),
	}).SignedString([]byte(os.Getenv("JWTKEY")))
	return tk, err
}

// Decode validates and decodes a raw JWT string into an Invite.
func DecodeInvite(raw string) (Invite, error) {
	if raw == "" {
		return Invite{}, ErrInvalid
	}

	claims := jwt.MapClaims{}
	if _, err := jwt.ParseWithClaims(
		raw, &claims, func(token *jwt.Token) (any, error) {
			return []byte(os.Getenv("JWTKEY")), nil
		},
	); err != nil {
		return Invite{}, err
	}

	return NewInvite(claims["teamID"].(string)), nil
}
