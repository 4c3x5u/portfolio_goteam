package login

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"server/assert"
	"server/auth"
	"server/db"
)

func TestHandler(t *testing.T) {
	var (
		readerPwd      = &db.FakeReaderUser{}
		comparerPwd    = &fakeComparer{}
		generatorToken = &auth.FakeGenerator{}
	)
	sut := NewHandler(readerPwd, comparerPwd, generatorToken)

	for _, c := range []struct {
		name                 string
		httpMethod           string
		reqBody              *ReqBody
		outResReaderUser     *db.User
		outErrReaderUser     error
		outResComparerHash   bool
		outErrComparerHash   error
		outResGeneratorToken string
		outErrGeneratorToken error
		wantStatusCode       int
	}{
		{
			name:                 "ErrHTTPMethod",
			httpMethod:           http.MethodGet,
			reqBody:              &ReqBody{},
			outResReaderUser:     nil,
			outErrReaderUser:     nil,
			outResComparerHash:   false,
			outErrComparerHash:   nil,
			outResGeneratorToken: "",
			outErrGeneratorToken: nil,
			wantStatusCode:       http.StatusMethodNotAllowed,
		},
		{
			name:                 "ErrNoUsername",
			httpMethod:           http.MethodPost,
			reqBody:              &ReqBody{},
			outResReaderUser:     nil,
			outErrReaderUser:     nil,
			outResComparerHash:   false,
			outErrComparerHash:   nil,
			outResGeneratorToken: "",
			outErrGeneratorToken: nil,
			wantStatusCode:       http.StatusBadRequest,
		},
		{
			name:                 "ErrUsernameEmpty",
			httpMethod:           http.MethodPost,
			reqBody:              &ReqBody{Username: ""},
			outResReaderUser:     nil,
			outErrReaderUser:     nil,
			outResComparerHash:   false,
			outErrComparerHash:   nil,
			outResGeneratorToken: "",
			outErrGeneratorToken: nil,
			wantStatusCode:       http.StatusBadRequest,
		},
		{
			name:                 "ErrUserNotFound",
			httpMethod:           http.MethodPost,
			reqBody:              &ReqBody{Username: "bob21"},
			outResReaderUser:     nil,
			outErrReaderUser:     sql.ErrNoRows,
			outResComparerHash:   false,
			outErrComparerHash:   nil,
			outResGeneratorToken: "",
			outErrGeneratorToken: nil,
			wantStatusCode:       http.StatusBadRequest,
		},
		{
			name:                 "ErrExistor",
			httpMethod:           http.MethodPost,
			reqBody:              &ReqBody{Username: "bob21", Password: "Myp4ssword!"},
			outResReaderUser:     nil,
			outErrReaderUser:     errors.New("existor fatal error"),
			outResComparerHash:   false,
			outErrComparerHash:   nil,
			outResGeneratorToken: "",
			outErrGeneratorToken: nil,
			wantStatusCode:       http.StatusInternalServerError,
		},
		{
			name:                 "ErrNoPassword",
			httpMethod:           http.MethodPost,
			reqBody:              &ReqBody{Username: "bob21"},
			outResReaderUser:     nil,
			outErrReaderUser:     nil,
			outResComparerHash:   false,
			outErrComparerHash:   nil,
			outResGeneratorToken: "",
			outErrGeneratorToken: nil,
			wantStatusCode:       http.StatusBadRequest,
		},
		{
			name:                 "ErrPasswordEmpty",
			httpMethod:           http.MethodPost,
			reqBody:              &ReqBody{Username: "bob21", Password: ""},
			outResReaderUser:     nil,
			outErrReaderUser:     nil,
			outResComparerHash:   false,
			outErrComparerHash:   nil,
			outResGeneratorToken: "",
			outErrGeneratorToken: nil,
			wantStatusCode:       http.StatusBadRequest,
		},
		{
			name:                 "ErrPasswordWrong",
			httpMethod:           http.MethodPost,
			reqBody:              &ReqBody{Username: "bob21", Password: "Myp4ssword!"},
			outResReaderUser:     &db.User{},
			outErrReaderUser:     nil,
			outResComparerHash:   false,
			outErrComparerHash:   nil,
			outResGeneratorToken: "",
			outErrGeneratorToken: nil,
			wantStatusCode:       http.StatusBadRequest,
		},
		{
			name:                 "ErrComparerHash",
			httpMethod:           http.MethodPost,
			reqBody:              &ReqBody{Username: "bob21", Password: "Myp4ssword!"},
			outResReaderUser:     &db.User{},
			outErrReaderUser:     nil,
			outResComparerHash:   true,
			outErrComparerHash:   errors.New("hash comparer error"),
			outResGeneratorToken: "",
			outErrGeneratorToken: nil,
			wantStatusCode:       http.StatusInternalServerError,
		},
		{
			name:                 "ErrGeneratorToken",
			httpMethod:           http.MethodPost,
			reqBody:              &ReqBody{Username: "bob21", Password: "Myp4ssword!"},
			outResReaderUser:     &db.User{},
			outErrReaderUser:     nil,
			outResComparerHash:   true,
			outErrComparerHash:   nil,
			outResGeneratorToken: "",
			outErrGeneratorToken: errors.New("token generator error"),
			wantStatusCode:       http.StatusInternalServerError,
		},
		{
			name:                 "OK",
			httpMethod:           http.MethodPost,
			reqBody:              &ReqBody{Username: "bob21", Password: "Myp4ssword!"},
			outResReaderUser:     &db.User{},
			outErrReaderUser:     nil,
			outResComparerHash:   true,
			outErrComparerHash:   nil,
			outResGeneratorToken: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9...",
			outErrGeneratorToken: nil,
			wantStatusCode:       http.StatusOK,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			readerPwd.OutRes = c.outResReaderUser
			readerPwd.OutErr = c.outErrReaderUser
			comparerPwd.outRes = c.outResComparerHash
			comparerPwd.outErr = c.outErrComparerHash
			generatorToken.OutRes = c.outResGeneratorToken
			generatorToken.OutErr = c.outErrGeneratorToken

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

			if err = assert.Equal(c.wantStatusCode, w.Result().StatusCode); err != nil {
				t.Error(err)
			}
			if c.wantStatusCode == http.StatusOK {
				// no field errors - assert on authToken
				foundSessionToken := false
				for _, cookie := range w.Result().Cookies() {
					if cookie.Name == "authToken" {
						foundSessionToken = true
						if err = assert.Equal(c.outResGeneratorToken, cookie.Value); err != nil {
							t.Error(err)
						}
					}
				}
				if err = assert.True(foundSessionToken); err != nil {
					t.Error(err)
				}
			}
		})
	}
}
