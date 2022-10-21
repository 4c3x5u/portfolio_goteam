package relay

import (
	"log"
	"net/http"
)

// APIMsger defines a type that uses a http response writer to relay messages.
type APIMsger interface {
	Msg(w http.ResponseWriter, msg string)
}

// APIErrMsger defines a type that uses a http response writer to relay messages
// and errors.
type APIErrMsger interface {
	Err(w http.ResponseWriter, errMsg string, status int)
	APIMsger
}

// APILogger is a means for the APILogger endpoints to log messages. It
// implements the APIErrMsger interface.
type APILogger struct {
}

// NewAPILogger is the constructor for APILogger. APILogger has no dependencies
// and this is written to emphasize the fact.
func NewAPILogger() *APILogger {
	return &APILogger{}
}

// Err relays an API error.
func (l *APILogger) Err(w http.ResponseWriter, errMsg string, status int) {
	http.Error(w, errMsg, status)
	log.Fatal(http.StatusText(status), "\n    ", errMsg, "\n")
}

// Msg relays an API message.
func (l *APILogger) Msg(w http.ResponseWriter, msg string) {
	log.Println(msg)
	if _, err := w.Write([]byte(msg)); err != nil {
		l.Err(w, err.Error(), http.StatusInternalServerError)
	}
}
