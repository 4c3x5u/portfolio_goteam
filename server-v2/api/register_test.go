package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/kxplxn/goteam/server-v2/relay"
)

func TestRegister(t *testing.T) {
	t.Run("Username Validation", func(t *testing.T) {
		t.Run("Too Short", func(t *testing.T) {
			// arrange
			req, err := http.NewRequest("POST", "/register", strings.NewReader(`{
				"username": "bob", 
				"password": "securepass1!", 
				"referrer": ""
			}`))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			handler := NewHandlerRegister(relay.NewAPILogger())
			wantErrs := []string{"Username cannot be shorter than 5 characters."}

			// act
			handler.ServeHTTP(w, req)

			// assert
			res := w.Result()
			if res.StatusCode != http.StatusBadRequest {
				t.Logf("\nwant: %d\ngot: %d", http.StatusBadRequest, res.StatusCode)
				t.Fail()
			}
			resBody := &ResRegister{}
			if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(resBody.Errs.Username, wantErrs) {
				t.Logf("\nwant: %+v\ngot: %+v", wantErrs, resBody.Errs.Username)
				t.Fail()
			}
		})
	})
}
