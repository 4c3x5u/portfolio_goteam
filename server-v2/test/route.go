// Package test contains code that is common between certain tests and
// encapsulates each type of test, test suite, and test case into types that
// contains the pieces of data needed to run the said common code.
package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kxplxn/goteam/server-v2/assert"
)

// RouteCase defines an API route test case.
type RouteCase struct {
	name     string
	input    string
	wantErrs []string
}

// NewRouteCase is the constructor for RouteCase.
func NewRouteCase(name string, input string, wantErrs []string) *RouteCase {
	return &RouteCase{name: name, input: input, wantErrs: wantErrs}
}

// RouteSuite contains a set of API route test cases as well as the name given
// to the set, the request/error field that the cases are concerned with, and
// the expected http response status code.
type RouteSuite struct {
	name           string
	field          string
	cases          []*RouteCase
	wantStatusCode int
}

// NewRouteSuite is the constructor for RouteSuite.
func NewRouteSuite(
	name string,
	inputField string,
	cases []*RouteCase,
	wantStatusCode int,
) *RouteSuite {
	return &RouteSuite{
		name:           name,
		field:          inputField,
		cases:          cases,
		wantStatusCode: wantStatusCode,
	}
}

// Route contains pieces of data that are common among a set of tests that are
// based on a request being made to a given API route.
type Route struct {
	url        string
	httpMethod string
	handler    http.Handler
	reqBody    map[string]string
	resBody    ErrsMapper
	suites     []*RouteSuite
}

// NewRoute is the constructor for Route.
func NewRoute(
	url string,
	httpMethod string,
	handler http.Handler,
	reqBody map[string]string,
	resBody ErrsMapper,
	suites []*RouteSuite,
) *Route {
	return &Route{
		url:        url,
		httpMethod: httpMethod,
		handler:    handler,
		reqBody:    reqBody,
		resBody:    resBody,
		suites:     suites,
	}
}

// Run runs a given set of RouteCase objects recursively.
func (r *Route) Run(t *testing.T) {
	for _, s := range r.suites {
		t.Run(s.name, func(t *testing.T) {
			for _, c := range s.cases {
				t.Run(c.name, func(t *testing.T) {
					// copy contents of valid request body example into a new map to avoid clashes (arrange)
					reqBody := make(map[string]string)
					for k, v := range r.reqBody {
						reqBody[k] = v
					}

					// set the specified request body field to the specified value (arrange)
					reqBody[s.field] = c.input
					reqBodyJSON, err := json.Marshal(reqBody)
					if err != nil {
						t.Fatal(err)
					}

					// set up the request (arrange)
					req, err := http.NewRequest(r.httpMethod, r.url, bytes.NewReader(reqBodyJSON))
					if err != nil {
						t.Fatal(err)
					}
					req.Header.Set("Content-Type", "application/json")
					w := httptest.NewRecorder()

					// send request (act)
					r.handler.ServeHTTP(w, req)

					// make assertions on the response, status code and errors returned (assert)
					res := w.Result()
					gotStatusCode := res.StatusCode
					if gotStatusCode != s.wantStatusCode {
						t.Logf("\nwant: %d\ngot: %d", s.wantStatusCode, res.StatusCode)
						t.Fail()
					}
					resBody := r.resBody
					if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
						t.Fatal(err)
					}
					gotErr := resBody.MapErrs()[s.field]
					if !assert.EqualArr(gotErr, c.wantErrs) {
						t.Logf("\nwant: %+v\ngot: %+v", c.wantErrs, gotErr)
						t.Fail()
					}
				})
			}
		})
	}
}
