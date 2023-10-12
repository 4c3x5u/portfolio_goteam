//go:build utest

package midware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kxplxn/goteam/server/assert"
)

// TestAccessControl tests the AccessControl middleware to assert that it
// returns the correct headers.
func TestAccessControl(t *testing.T) {
	wantAllowOrigin := "http://allowedorigin.com"
	wantAllowHeaders := "Content-Type"
	wantAllowCredentials := "true"

	sut := NewAccessControl(nil, wantAllowOrigin)
	req, err := http.NewRequest(http.MethodOptions, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	sut.ServeHTTP(w, req)

	if err = assert.Equal(http.StatusOK, w.Result().StatusCode); err != nil {
		t.Error(err)
	}
	if err = assert.Equal(
		wantAllowOrigin,
		w.Result().Header.Get("Access-Control-Allow-Origin"),
	); err != nil {
		t.Error(err)
	}
	if err = assert.Equal(
		wantAllowHeaders,
		w.Result().Header.Get("Access-Control-Allow-Headers"),
	); err != nil {
		t.Error(err)
	}
	if err = assert.Equal(
		wantAllowCredentials,
		w.Result().Header.Get("Access-Control-Allow-Credentials"),
	); err != nil {
		t.Error(err)
	}
}
