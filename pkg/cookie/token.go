// Package token contains code for generating, validating, and decoding JWTs.
package cookie

import (
	"errors"
	"net/http"
)

// Encoder defines a type that can be used to encode a JWT.
type Encoder[T any] interface{ Encode(T) (http.Cookie, error) }

// Decoder defines a type that can be used to decode a JWT.
type Decoder[T any] interface{ Decode(http.Cookie) (T, error) }

// ErrInvalid means that the given token was invalid.
var ErrInvalid = errors.New("invalid token")
