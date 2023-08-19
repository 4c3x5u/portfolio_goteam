//go:build utest

package midware

import (
	"net/http"
	"net/http/httptest"
	"server/assert"
	"testing"
)

// TestCORS tests the CORS middleware to assert that it returns the correct
// headers.
func TestCORS(t *testing.T) {
	wantAllowedOrigin := "http://allowedorigin.com"
	wantAllowedHeaders := "Content-Type"
	sut := NewCORS(nil, wantAllowedOrigin)
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
		wantAllowedOrigin,
		w.Result().Header.Get("Access-Control-Allow-Origin"),
	); err != nil {
		t.Error(err)
	}
	if err = assert.Equal(
		wantAllowedHeaders,
		w.Result().Header.Get("Access-Control-Allow-Headers"),
	); err != nil {
		t.Error(err)
	}
}
