package login

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"server/assert"
	"server/db"
)

func TestHandler(t *testing.T) {
	existorUser := &db.FakeExistor{}
	sut := NewHandler(existorUser)

	t.Run("ErrHttpMethod", func(t *testing.T) {
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

	t.Run("ErrNoUsername", func(t *testing.T) {
		reqBody := &ReqBody{}
		reqBodyJSON, err := json.Marshal(reqBody)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqBodyJSON))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		sut.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})
}
