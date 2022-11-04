package relay

import (
	"log"
	"net/http"
)

// ErrAPIInternal relays an internal API error by logging a message to the
// console and returning 500 to the user.
func ErrAPIInternal(w http.ResponseWriter, msg string) {
	log.Printf("ERROR: %s", msg)
	w.WriteHeader(http.StatusInternalServerError)
	return
}
