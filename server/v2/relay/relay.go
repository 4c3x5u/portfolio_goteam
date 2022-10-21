// Package relay contains functions used by the server to communicate with the
// various other parts of the system (terminal, web console, etc.).
package relay

import (
	"log"
	"net/http"
)

// ErrMsger defines a type that can be initialised with a http response
// writer and relay messages and errors.
type ErrMsger interface {
	Err(w http.ResponseWriter, errMsg string, status int)
	Msg(w http.ResponseWriter, msg string)
}

// Logger is a means for the API endpoints to log messages.
type Logger struct {
}

// NewLogger is the constructor for Logger. Logger has no dependencies
// and this is written to emphasize the fact.
func NewLogger() *Logger {
	return &Logger{}
}

// Err relays an API error.
func (a *Logger) Err(w http.ResponseWriter, errMsg string, status int) {
	http.Error(w, errMsg, status)
	log.Fatal(http.StatusText(status), "\n    ", errMsg, "\n")
}

// Msg relays an API message.
func (a *Logger) Msg(w http.ResponseWriter, msg string) {
	log.Println(msg)
	if _, err := w.Write([]byte(msg)); err != nil {
		a.Err(w, err.Error(), http.StatusInternalServerError)
	}
}
