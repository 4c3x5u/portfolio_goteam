// go:build itest

// Package itest contains integration tests for the application.
package itest

// serverHost is the host that the server runs on in the test container.
const serverHost = "localhost"

// serverPort is the port that the server runs at in test container.
const serverPort = "8081"

// serverURL is the url that is used to send requests to the server running in
// the test container. It is set during setup in main_test.go/MainTest.
var serverURL = "http://" + serverHost + ":" + serverPort
