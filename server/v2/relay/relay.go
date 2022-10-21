// Package relay (verb) contains functions used by the web api to communicate
// with the various parts of the system.
package relay

import (
	"log"
	"net/http"
)

type APIErrMsger interface {
	APIErr(w http.ResponseWriter, errMsg string, status int)
	APIMsg(w http.ResponseWriter, msg string)
}

type Relay struct {
	w http.ResponseWriter
}

func New(w http.ResponseWriter) *Relay {
	return &Relay{w}
}

func (r *Relay) APIMsg(msg string) {
	log.Println(msg)
	if _, err := r.w.Write([]byte(msg)); err != nil {
		r.APIErr(err.Error(), http.StatusInternalServerError)
	}
}

// APIErr relays an error.
func (r *Relay) APIErr(errMsg string, status int) {
	http.Error(r.w, errMsg, status)
	log.Fatal(http.StatusText(status), "\n    ", errMsg, "\n")
}
