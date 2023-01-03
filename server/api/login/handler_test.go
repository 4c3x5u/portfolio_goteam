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
	"server/db"
)

func TestHandler(t *testing.T) {
	var (
		readerPwd     = &db.FakeReaderUser{}
		comparerPwd   = &fakeComparer{}
		readerSession = &db.FakeReaderSession{}
	)
	sut := NewHandler(readerPwd, comparerPwd, readerSession)

	for _, c := range []struct {
		name                string
		httpMethod          string
		reqBody             *ReqBody
		outResReaderUser    *db.User
		outErrReaderUser    error
		outResComparerHash  bool
		outErrComparerHash  error
		outErrReaderSession error
		wantStatusCode      int
	}{
		{
			name:                "ErrHTTPMethod",
			httpMethod:          http.MethodGet,
			reqBody:             &ReqBody{},
			outResReaderUser:    nil,
			outErrReaderUser:    nil,
			outResComparerHash:  false,
			outErrComparerHash:  nil,
			outErrReaderSession: nil,
			wantStatusCode:      http.StatusMethodNotAllowed,
		},
		{
			name:                "ErrNoUsername",
			httpMethod:          http.MethodPost,
			reqBody:             &ReqBody{},
			outResReaderUser:    nil,
			outErrReaderUser:    nil,
			outResComparerHash:  false,
			outErrComparerHash:  nil,
			outErrReaderSession: nil,
			wantStatusCode:      http.StatusBadRequest,
		},
		{
			name:                "ErrUsernameEmpty",
			httpMethod:          http.MethodPost,
			reqBody:             &ReqBody{Username: ""},
			outResReaderUser:    nil,
			outErrReaderUser:    nil,
			outResComparerHash:  false,
			outErrComparerHash:  nil,
			outErrReaderSession: nil,
			wantStatusCode:      http.StatusBadRequest,
		},
		{
			name:                "ErrUserNotFound",
			httpMethod:          http.MethodPost,
			reqBody:             &ReqBody{Username: "bob21"},
			outResReaderUser:    nil,
			outErrReaderUser:    sql.ErrNoRows,
			outResComparerHash:  false,
			outErrComparerHash:  nil,
			outErrReaderSession: nil,
			wantStatusCode:      http.StatusBadRequest,
		},
		{
			name:                "ErrExistor",
			httpMethod:          http.MethodPost,
			reqBody:             &ReqBody{Username: "bob21", Password: "Myp4ssword!"},
			outResReaderUser:    nil,
			outErrReaderUser:    errors.New("existor fatal error"),
			outResComparerHash:  false,
			outErrComparerHash:  nil,
			outErrReaderSession: nil,
			wantStatusCode:      http.StatusInternalServerError,
		},
		{
			name:                "ErrNoPassword",
			httpMethod:          http.MethodPost,
			reqBody:             &ReqBody{Username: "bob21"},
			outResReaderUser:    nil,
			outErrReaderUser:    nil,
			outResComparerHash:  false,
			outErrComparerHash:  nil,
			outErrReaderSession: nil,
			wantStatusCode:      http.StatusBadRequest,
		},
		{
			name:                "ErrPasswordEmpty",
			httpMethod:          http.MethodPost,
			reqBody:             &ReqBody{Username: "bob21", Password: ""},
			outResReaderUser:    nil,
			outErrReaderUser:    nil,
			outResComparerHash:  false,
			outErrComparerHash:  nil,
			outErrReaderSession: nil,
			wantStatusCode:      http.StatusBadRequest,
		},
		{
			name:                "ErrPasswordWrong",
			httpMethod:          http.MethodPost,
			reqBody:             &ReqBody{Username: "bob21", Password: "Myp4ssword!"},
			outResReaderUser:    &db.User{},
			outErrReaderUser:    nil,
			outResComparerHash:  false,
			outErrComparerHash:  nil,
			outErrReaderSession: nil,
			wantStatusCode:      http.StatusBadRequest,
		},
		{
			name:                "ErrComparerHash",
			httpMethod:          http.MethodPost,
			reqBody:             &ReqBody{Username: "bob21", Password: "Myp4ssword!"},
			outResReaderUser:    &db.User{},
			outErrReaderUser:    nil,
			outResComparerHash:  true,
			outErrComparerHash:  nil,
			outErrReaderSession: errors.New("session reader error"),
			wantStatusCode:      http.StatusInternalServerError,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			readerPwd.OutRes = c.outResReaderUser
			readerPwd.OutErr = c.outErrReaderUser
			comparerPwd.outRes = c.outResComparerHash
			comparerPwd.outErr = c.outErrComparerHash
			readerSession.OutErr = c.outErrReaderSession

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
