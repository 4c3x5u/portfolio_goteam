package cookie

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// InviteName is the name of the invite token.
const InviteName = "invite-token"

// Invite defines the body of an Invite token.
type Invite struct{ TeamID string }

// NewInvite creates and returns a new Invite.
func NewInvite(teamID string) Invite { return Invite{TeamID: teamID} }

// InviteEncoder defines a type that can be used to encode an invite token.
type InviteEncoder struct {
	key []byte
	dur time.Duration
}

// NewInviteEncoder creates and returns a new InviteEncoder.
func NewInviteEncoder(key []byte, dur time.Duration) InviteEncoder {
	return InviteEncoder{key: key, dur: dur}
}

// Encode encodes an Invite into a JWT string.
func (e InviteEncoder) Encode(inv Invite) (http.Cookie, error) {
	exp := time.Now().Add(e.dur)

	tk, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"teamID": inv.TeamID,
		"exp":    exp.Unix(),
	}).SignedString(e.key)
	if err != nil {
		return http.Cookie{}, err
	}

	return http.Cookie{
		Name:     InviteName,
		Value:    tk,
		Expires:  exp.UTC(),
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}, nil
}

// InviteDecoder defines a type that can be used to decode an invite token.
type InviteDecoder struct{ key []byte }

// NewInviteDecoder creates and returns a new InviteDecoder.
func NewInviteDecoder(key []byte) InviteDecoder {
	return InviteDecoder{key: key}
}

// Decode validates and decodes a raw JWT string into an Invite.
func (d InviteDecoder) Decode(token string) (Invite, error) {
	claims := jwt.MapClaims{}
	if _, err := jwt.ParseWithClaims(
		token, &claims, func(token *jwt.Token) (any, error) {
			return d.key, nil
		},
	); err != nil {
		return Invite{}, err
	}

	return NewInvite(claims["teamID"].(string)), nil
}
