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
	Err(w http.ResponseWriter, statusCode int)
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

// Msg relays an API message.
func (l *APILogger) Msg(w http.ResponseWriter, msg string) {
	if _, err := w.Write([]byte(msg)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err.Error())
	}
	log.Println(msg)
}

// Err relays an error based on a HTTP status code.
func (l *APILogger) Err(w http.ResponseWriter, statusCode int) {
	http.Error(w, http.StatusText(statusCode), statusCode)
	log.Fatal(http.StatusText(statusCode))
}
