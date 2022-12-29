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

	for _, c := range []struct {
		name           string
		httpMethod     string
		reqBody        *ReqBody
		wantStatusCode int
	}{
		{
			name:           "ErrHTTPMethod",
			httpMethod:     http.MethodGet,
			reqBody:        &ReqBody{},
			wantStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:           "ErrNoUsername",
			httpMethod:     http.MethodPost,
			reqBody:        &ReqBody{},
			wantStatusCode: http.StatusBadRequest,
		},
	} {
		reqBodyJSON, err := json.Marshal(c.reqBody)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest(c.httpMethod, "/login", bytes.NewReader(reqBodyJSON))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		sut.ServeHTTP(w, req)

		assert.Equal(t, c.wantStatusCode, w.Result().StatusCode)
	}
}
