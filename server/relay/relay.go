// Package relay contains code used to relay messages to either the end user,
// the, logs or both.
//
// TODO: add tests
package relay

import (
	"encoding/json"
	"log"
	"net/http"
)

// ServerErr relays an internal server error by logging a message to the console
// and returning 500 to the user.
func ServerErr(w http.ResponseWriter, msg string) {
	log.Printf("ERROR: %s", msg)
	w.WriteHeader(http.StatusInternalServerError)
}

// ClientJSON relays a JSON object to the client by writing it to the HTTP
// response body as well as writing the specified status code into the header.
func ClientJSON(w http.ResponseWriter, res any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		ServerErr(w, err.Error())
	}
}
