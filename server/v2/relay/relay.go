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
	ErrCode(w http.ResponseWriter, statusCode int)
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

// ErrCode relays an error based on a HTTP status code.
func (l *APILogger) ErrCode(w http.ResponseWriter, statusCode int) {
	msg := http.StatusText(statusCode)
	http.Error(w, msg, statusCode)
	log.Fatal(http.StatusText(statusCode), "\n    ", msg, "\n")
}
