//go:build utest

package api

import (
	"net/http"
	"server/assert"
	"testing"
)

// TestAllowedMethods tests the AllowedMethods function to assert that it
// returns the correct header alongside the comma-separated list of strings
// passed to it as a string.
func TestAllowedMethods(t *testing.T) {
	methods := []string{http.MethodGet, http.MethodDelete, http.MethodPatch}
	header, allowedMethods := AllowedMethods(methods)
	if err := assert.Equal("Access-Control-Allow-Methods", header); err != nil {
		t.Error(err)
	}
	if err := assert.Equal("GET, DELETE, PATCH", allowedMethods); err != nil {
		t.Error(err)
	}
}
