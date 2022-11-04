package api

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kxplxn/goteam/server-v2/relay"
)

func TestRegister(t *testing.T) {
	t.Run("Username Validation", func(t *testing.T) {
		t.Run("Too Short", func(t *testing.T) {
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
			wantRes := &ResRegister{
				Errs: &ErrsRegister{
					Username: []string{"Username cannot be shorter than 5 characters."},
				},
			}

			handler.ServeHTTP(w, req)

			httpResponse := w.Result()
			if httpResponse.StatusCode != http.StatusBadRequest {
				t.Logf("\nwant: %d\ngot: %d", http.StatusBadRequest, httpResponse.StatusCode)
				t.Fail()
			}
			gotRes := &ResRegister{Errs: &ErrsRegister{}}
			if err := json.NewDecoder(httpResponse.Body).Decode(&gotRes); err != nil {
				t.Fatal(err)
			}
			if cmp.Equal(gotRes, wantRes) == false {
				t.Logf("\nwant: %v\ngot: %v", wantRes, gotRes)
				t.Fail()
			}
		})
	})
}
