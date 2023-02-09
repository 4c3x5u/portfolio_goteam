//go:build itest

package itest

import (
	"errors"
	"net/http"
	"testing"
)

// TestBoardRoute tests the /board route to assert that it behaves correctly.
func TestBoardRoute(t *testing.T) {
	t.Run("healthcheck", func(t *testing.T) {
		if res, err := http.Get(serverURL); err != nil {
			t.Fatal(err)
		} else if res.StatusCode != 200 {
			t.Fatal(errors.New("status: " + res.Status))
		}
	})

	// TODO: test the board route.
}
