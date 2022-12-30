package login

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"server/db"
	"testing"

	"server/assert"
)

func TestHandler(t *testing.T) {
	existorUser := &db.FakeExistor{}
	sut := NewHandler(existorUser)

	for _, c := range []struct {
		name              string
		httpMethod        string
		reqBody           *ReqBody
		wantStatusCode    int
		outResExistorUser bool
	}{
		{
			name:              "ErrHTTPMethod",
			httpMethod:        http.MethodGet,
			reqBody:           &ReqBody{},
			wantStatusCode:    http.StatusMethodNotAllowed,
			outResExistorUser: true,
		},
		{
			name:              "ErrNoUsername",
			httpMethod:        http.MethodPost,
			reqBody:           &ReqBody{},
			wantStatusCode:    http.StatusBadRequest,
			outResExistorUser: true,
		},
		{
			name:              "ErrUsernameEmpty",
			httpMethod:        http.MethodPost,
			reqBody:           &ReqBody{Username: ""},
			wantStatusCode:    http.StatusBadRequest,
			outResExistorUser: true,
		},
		{
			name:              "ErrUserNotFound",
			httpMethod:        http.MethodPost,
			reqBody:           &ReqBody{Username: "bob21"},
			wantStatusCode:    http.StatusBadRequest,
			outResExistorUser: false,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			existorUser.OutExists = c.outResExistorUser

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
		})
	}
}
