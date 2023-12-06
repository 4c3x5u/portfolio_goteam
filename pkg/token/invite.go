package token

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// NameInvite is the name of the invite token.
const NameInvite = "invite"

// Invite defines the body of an invite token.
type Invite struct{ TeamID string }

// Encode encodes the Invite into a JWT string
func (i *Invite) Encode(exp time.Time) (string, error) {
	tk, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"teamID": i.TeamID,
		"exp":    exp.Unix(),
	}).SignedString([]byte(signKey))
	return tk, err
}

// Decode validates and decodes a raw JWT string into the Invite.
func (i *Invite) Decode(raw string) error {
	if raw == "" {
		return ErrInvalid
	}

	claims := jwt.MapClaims{}
	if _, err := jwt.ParseWithClaims(
		raw, &claims, func(token *jwt.Token) (any, error) {
			return signKey, nil
		},
	); err != nil {
		return err
	}

	i.TeamID = claims["teamID"].(string)

	return nil
}
