// Package relay contains functions used by the web api to communicate with the
// various parts of the system.
package relay

import (
	"log"
	"net/http"
)

// APIErrMsger defines a type that can be initialised with a http response
// writer and relay messages and errors.
type APIErrMsger interface {
	Err(w http.ResponseWriter, errMsg string, status int)
	Msg(w http.ResponseWriter, msg string)
}

// APILogger is a means for the API endpoints to log messages.
type APILogger struct {
}

// NewAPILogger is the constructor for APILogger. APILogger has no dependencies
// and this is written to emphasize the fact.
func NewAPILogger() *APILogger {
	return &APILogger{}
}

// Err relays an API error.
func (a *APILogger) Err(w http.ResponseWriter, errMsg string, status int) {
	http.Error(w, errMsg, status)
	log.Fatal(http.StatusText(status), "\n    ", errMsg, "\n")
}

// Msg relays an API message.
func (a *APILogger) Msg(w http.ResponseWriter, msg string) {
	log.Println(msg)
	if _, err := w.Write([]byte(msg)); err != nil {
		a.Err(w, err.Error(), http.StatusInternalServerError)
	}
}
