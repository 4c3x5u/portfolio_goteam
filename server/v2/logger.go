package main

import (
	"log"
	"net/http"
)

type Logger struct {
	w http.ResponseWriter
}

func NewLogger(w http.ResponseWriter) *Logger {
	return &Logger{w: w}
}

func (l *Logger) Msg(msg string) {
	log.Println(msg)
	if _, err := l.w.Write([]byte(msg)); err != nil {
		l.Err(err.Error(), http.StatusInternalServerError)
	}
}

func (l *Logger) Err(errMsg string, status int) {
	http.Error(l.w, errMsg, status)
	log.Fatal(http.StatusText(status), "\n    ", errMsg)
}
