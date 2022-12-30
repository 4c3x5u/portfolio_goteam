package login

import (
	"bytes"
	"encoding/json"
	"errors"
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
		name              string
		httpMethod        string
		reqBody           *ReqBody
		outResExistorUser bool
		outErrExistorUser error
		wantStatusCode    int
	}{
		{
			name:              "ErrHTTPMethod",
			httpMethod:        http.MethodGet,
			reqBody:           &ReqBody{},
			outResExistorUser: true,
			outErrExistorUser: nil,
			wantStatusCode:    http.StatusMethodNotAllowed,
		},
		{
			name:              "ErrNoUsername",
			httpMethod:        http.MethodPost,
			reqBody:           &ReqBody{},
			outResExistorUser: true,
			outErrExistorUser: nil,
			wantStatusCode:    http.StatusBadRequest,
		},
		{
			name:              "ErrUsernameEmpty",
			httpMethod:        http.MethodPost,
			reqBody:           &ReqBody{Username: ""},
			outResExistorUser: true,
			outErrExistorUser: nil,
			wantStatusCode:    http.StatusBadRequest,
		},
		{
			name:              "ErrUserNotFound",
			httpMethod:        http.MethodPost,
			reqBody:           &ReqBody{Username: "bob21"},
			outResExistorUser: false,
			outErrExistorUser: nil,
			wantStatusCode:    http.StatusBadRequest,
		},
		{
			name:              "ErrExistor",
			httpMethod:        http.MethodPost,
			reqBody:           &ReqBody{Username: "bob21", Password: "Myp4ssword!"},
			outResExistorUser: true,
			outErrExistorUser: errors.New("existor fatal error"),
			wantStatusCode:    http.StatusInternalServerError,
		},
		{
			name:              "ErrNoPassword",
			httpMethod:        http.MethodPost,
			reqBody:           &ReqBody{Username: "bob21"},
			outResExistorUser: true,
			outErrExistorUser: nil,
			wantStatusCode:    http.StatusBadRequest,
		},
		{
			name:              "ErrPasswordEmpty",
			httpMethod:        http.MethodPost,
			reqBody:           &ReqBody{Username: "bob21", Password: ""},
			outResExistorUser: true,
			outErrExistorUser: nil,
			wantStatusCode:    http.StatusBadRequest,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			existorUser.OutExists = c.outResExistorUser
			existorUser.OutErr = c.outErrExistorUser

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
