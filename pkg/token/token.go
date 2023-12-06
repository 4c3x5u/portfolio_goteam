// Package token contains code for generating, validating, and decoding JWTs.
package token

import (
	"errors"
	"os"
	"time"
)

// signKey can be used to sign and parse JWTs.
var signKey []byte

func init() { signKey = []byte(os.Getenv("JWTKEY")) }

// ErrInvalid means that the given token is invalid.
var ErrInvalid = errors.New("invalid token")

// Encoder defines a type that can be used to encode a token.
type EncodeFunc[T any] func(T, time.Time) (string, error)

// Decoder defines a type that can be used to decode a JWT.
type Decoder[T any] func(string) (T, error)
