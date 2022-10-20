// Package relay (verb) contains functions used by the web api to communicate
// with the various parts of the system.
package relay

import (
	"log"
	"net/http"
)

// APIMsg relays an API message.
func APIMsg(w http.ResponseWriter, msg string) {
	log.Println(msg)
	if _, err := w.Write([]byte(msg)); err != nil {
		APIErr(w, err.Error(), http.StatusInternalServerError)
	}
}

// APIErr relays an error.
func APIErr(w http.ResponseWriter, errMsg string, status int) {
	http.Error(w, errMsg, status)
	log.Fatal(http.StatusText(status), "\n    ", errMsg, "\n")
}
