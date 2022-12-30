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
	readerPwd := &db.FakeReaderBytes{}
	comparerPwd := &fakeComparer{}
	sut := NewHandler(readerPwd, comparerPwd)

	for _, c := range []struct {
		name                string
		httpMethod          string
		reqBody             *ReqBody
		outResReaderUserPwd []byte
		outErrReaderUserPwd error
		outResComparerHash  bool
		outErrComparerHash  error
		wantStatusCode      int
	}{
		{
			name:                "ErrHTTPMethod",
			httpMethod:          http.MethodGet,
			reqBody:             &ReqBody{},
			outResReaderUserPwd: []byte{},
			outErrReaderUserPwd: nil,
			outResComparerHash:  false,
			outErrComparerHash:  nil,
			wantStatusCode:      http.StatusMethodNotAllowed,
		},
		{
			name:                "ErrNoUsername",
			httpMethod:          http.MethodPost,
			reqBody:             &ReqBody{},
			outResReaderUserPwd: []byte{},
			outErrReaderUserPwd: nil,
			outResComparerHash:  false,
			outErrComparerHash:  nil,
			wantStatusCode:      http.StatusBadRequest,
		},
		{
			name:                "ErrUsernameEmpty",
			httpMethod:          http.MethodPost,
			reqBody:             &ReqBody{Username: ""},
			outResReaderUserPwd: []byte{},
			outErrReaderUserPwd: nil,
			outResComparerHash:  false,
			outErrComparerHash:  nil,
			wantStatusCode:      http.StatusBadRequest,
		},
		{
			name:                "ErrUserNotFound",
			httpMethod:          http.MethodPost,
			reqBody:             &ReqBody{Username: "bob21"},
			outResReaderUserPwd: []byte{},
			outErrReaderUserPwd: sql.ErrNoRows,
			outResComparerHash:  false,
			outErrComparerHash:  nil,
			wantStatusCode:      http.StatusBadRequest,
		},
		{
			name:                "ErrExistor",
			httpMethod:          http.MethodPost,
			reqBody:             &ReqBody{Username: "bob21", Password: "Myp4ssword!"},
			outResReaderUserPwd: []byte{},
			outErrReaderUserPwd: errors.New("existor fatal error"),
			outResComparerHash:  false,
			outErrComparerHash:  nil,
			wantStatusCode:      http.StatusInternalServerError,
		},
		{
			name:                "ErrNoPassword",
			httpMethod:          http.MethodPost,
			reqBody:             &ReqBody{Username: "bob21"},
			outResReaderUserPwd: []byte{},
			outErrReaderUserPwd: nil,
			outResComparerHash:  false,
			outErrComparerHash:  nil,
			wantStatusCode:      http.StatusBadRequest,
		},
		{
			name:                "ErrPasswordEmpty",
			httpMethod:          http.MethodPost,
			reqBody:             &ReqBody{Username: "bob21", Password: ""},
			outResReaderUserPwd: []byte{},
			outErrReaderUserPwd: nil,
			outResComparerHash:  false,
			outErrComparerHash:  nil,
			wantStatusCode:      http.StatusBadRequest,
		},
		{
			name:                "ErrPasswordWrong",
			httpMethod:          http.MethodPost,
			reqBody:             &ReqBody{Username: "bob21", Password: "Myp4ssword!"},
			outResReaderUserPwd: []byte{},
			outErrReaderUserPwd: nil,
			outResComparerHash:  false,
			outErrComparerHash:  nil,
			wantStatusCode:      http.StatusBadRequest,
		},
		{
			name:                "ErrComparerHash",
			httpMethod:          http.MethodPost,
			reqBody:             &ReqBody{Username: "bob21", Password: "Myp4ssword!"},
			outResReaderUserPwd: []byte{},
			outErrReaderUserPwd: nil,
			outResComparerHash:  false,
			outErrComparerHash:  errors.New("comparer fatal error"),
			wantStatusCode:      http.StatusInternalServerError,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			readerPwd.OutRes = c.outResReaderUserPwd
			readerPwd.OutErr = c.outErrReaderUserPwd
			comparerPwd.outRes = c.outResComparerHash
			comparerPwd.outErr = c.outErrComparerHash

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
