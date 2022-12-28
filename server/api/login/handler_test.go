package login

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"server/assert"
)

func TestHandler(t *testing.T) {
	t.Run("ErrHttpMethod", func(t *testing.T) {
		sut := NewHandler()

		reqBody, err := json.Marshal(&ReqBody{})
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest(http.MethodGet, "/login", bytes.NewReader(reqBody))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		sut.ServeHTTP(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Result().StatusCode)
	})
}
