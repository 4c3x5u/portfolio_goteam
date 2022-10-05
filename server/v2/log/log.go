package log

import (
	"log"
	"net/http"
)

func ErrToConsole(err error) {
	log.Fatal(err.Error())
}

type APILogger struct {
	w http.ResponseWriter
}

func NewAPILogger(w http.ResponseWriter) *APILogger {
	return &APILogger{w: w}
}

func (l *APILogger) Msg(msg string) {
	log.Println(msg)
	if _, err := l.w.Write([]byte(msg)); err != nil {
		l.Err(err.Error(), http.StatusInternalServerError)
	}
}

func (l *APILogger) Err(errMsg string, status int) {
	http.Error(l.w, errMsg, status)
	log.Fatal(http.StatusText(status), "\n    ", errMsg)
}
