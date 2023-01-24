package board

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"server/assert"
)

// TestDELETEHandler tests the Handle method of DELETEHandler to assert that
// it behaves correctly.
func TestDELETEHandler(t *testing.T) {
	sut := NewDELETEHandler()
	req, err := http.NewRequest(http.MethodPost, "/board?id=123", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	sut.Handle(w, req, "bob123")

	if err := assert.Equal(
		http.StatusNotImplemented, w.Result().StatusCode,
	); err != nil {
		t.Error(err)
	}
}
