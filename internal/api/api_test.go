//go:build utest

package api

import (
	"net/http"
	"testing"

	"github.com/kxplxn/goteam/pkg/assert"
)

// TestAllowedMethods tests the AllowedMethods function to assert that it
// returns the correct header alongside the comma-separated list of strings
// passed to it as a string.
func TestAllowedMethods(t *testing.T) {
	methods := []string{http.MethodGet, http.MethodDelete, http.MethodPatch}
	header, allowedMethods := AllowedMethods(methods)
	assert.Equal(t.Error, header, "Access-Control-Allow-Methods")
	assert.Equal(t.Error, allowedMethods, "GET, DELETE, PATCH")
}
