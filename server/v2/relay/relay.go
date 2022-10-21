package relay

import (
	"log"
	"net/http"
)

// APIMsger defines a type that uses a http response writer to relay messages.
type APIMsger interface {
	Msg(w http.ResponseWriter, msg string)
}

// APIErrMsger defines a type that uses a http response writer to relay
// errors based on a HTTP status code, as well as relaying messages.
type APIErrMsger interface {
	ErrStatus(w http.ResponseWriter, statusCode int)
	Err(w http.ResponseWriter, msg string, statusCode int)
	APIMsger
}

// APILogger is a means for the APILogger endpoints to log messages. It
// implements APIErrMsger.
type APILogger struct {
}

// NewAPILogger is the constructor for APILogger. APILogger has no dependencies
// and this is written to emphasize the fact.
func NewAPILogger() *APILogger {
	return &APILogger{}
}

// Err relays an API error using a message text and HTTP status code.
func (l *APILogger) Err(w http.ResponseWriter, msg string, statusCode int) {
	http.Error(w, msg, statusCode)
	log.Fatal(msg)
}

// ErrStatus relays an API error based on a HTTP status code.
func (l *APILogger) ErrStatus(w http.ResponseWriter, statusCode int) {
	l.Err(w, http.StatusText(statusCode), statusCode)
}

// Msg relays an API message.
func (l *APILogger) Msg(w http.ResponseWriter, msg string) {
	if _, err := w.Write([]byte(msg)); err != nil {
		l.Err(w, err.Error(), http.StatusInternalServerError)
	}
	log.Println(msg)
}
