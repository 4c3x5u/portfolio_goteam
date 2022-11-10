package test

import (
	"bytes"
	"encoding/json"
	"github.com/kxplxn/goteam/server-v2/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ErrsMapper allows errors returned from API routes to be investigated in
// Route.Run through the field name that is set on each RoutCase.
type ErrsMapper interface {
	MapErrs() map[string][]string
}

// RouteCase defines an API route test case.
type RouteCase struct {
	name       string
	inputField string
	input      string
	errsField  string
	wantErrs   []string
	children   []*RouteCase
}

// Route contains pieces of data that are common among a set of RouteCase
// objects. It does NOT contain the RouteCase objects themselves for its Run
// function to be able to be called recursively.
type Route struct {
	address      string
	httpMethod   string
	handler      http.Handler
	validReqBody map[string]string
	resBody      ErrsMapper
}

// Run runs a given set of RouteCase objects recursively.
func (r *Route) Run(t *testing.T, cases []*RouteCase) {
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.children != nil && len(c.children) > 0 {
				r.Run(t, c.children)
			}

			// arrange
			reqBody := r.validReqBody
			reqBody[c.inputField] = c.input
			reqBodyJSON, err := json.Marshal(reqBody)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(r.httpMethod, r.address, bytes.NewReader(reqBodyJSON))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// act
			r.handler.ServeHTTP(w, req)

			// assert
			res := w.Result()
			gotStatusCode, wantStatusCode := res.StatusCode, http.StatusBadRequest
			if gotStatusCode != wantStatusCode {
				t.Logf("\nwant: %d\ngot: %d", http.StatusBadRequest, res.StatusCode)
				t.Fail()
			}
			if err := json.NewDecoder(res.Body).Decode(&r.resBody); err != nil {
				t.Fatal(err)
			}
			gotErr := r.resBody.MapErrs()[c.errsField]
			if !assert.EqualArr(gotErr, c.wantErrs) {
				t.Logf("\nwant: %+v\ngot: %+v", c.wantErrs, gotErr)
				t.Fail()
			}
		})
	}
}
