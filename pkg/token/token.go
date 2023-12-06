// Package token contains code for generating, validating, and decoding JWTs.
package token

import (
	"errors"
	"os"
)

// signKey can be used to sign and parse JWTs.
var signKey []byte

func init() { signKey = []byte(os.Getenv("JWTKEY")) }

// ErrInvalid means that the given token is invalid.
var ErrInvalid = errors.New("invalid token")

// DecodeFunc defines a function that can be used to decode a token.
type DecodeFunc[T any] func(string) (T, error)
