// Package itest contains integration tests for the application. Since it is an
// actual module it will only run when `go test` is executed for it explicitly,
// unlike the rest of the tests in other packages which can run by executing
// `go test run ./...` from server directory.
//
// Each of the test files correspond to an API route that the app serves, except
// the main_test.go file which handles setup for the integration tests and runs
// them.
package itest
