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

// APILogger is a means for the APILogger endpoints to log messages.
type APILogger struct {
}

// NewAPILogger is the constructor for APILogger. APILogger has no dependencies
// and this is written to emphasize the fact.
func NewAPILogger() *APILogger {
	return &APILogger{}
}

// Err relays an APILogger error.
func (a *APILogger) Err(w http.ResponseWriter, errMsg string, status int) {
	http.Error(w, errMsg, status)
	log.Fatal(http.StatusText(status), "\n    ", errMsg, "\n")
}

// Msg relays an APILogger message.
func (a *APILogger) Msg(w http.ResponseWriter, msg string) {
	log.Println(msg)
	if _, err := w.Write([]byte(msg)); err != nil {
		a.Err(w, err.Error(), http.StatusInternalServerError)
	}
}
