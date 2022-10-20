package api

import (
	"encoding/json"
	"fmt"
	"github.com/kxplxn/goteam/server/v2/relay"
	"net/http"
)

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	// accept only POST
	if r.Method != "POST" {
		status := http.StatusMethodNotAllowed
		relay.APIErr(w, http.StatusText(status), status)
		return
	}

	// read body into map
	dec := make(map[string]string, 3)
	if err := json.NewDecoder(r.Body).Decode(&dec); err != nil {
		status := http.StatusInternalServerError
		relay.APIErr(w, http.StatusText(status), status)
	}

	// log decoded body
	relay.APIMsg(w, fmt.Sprintf(
		"usn: %s\npwd: %s\nref: %s\n",
		dec["usn"], dec["pwd"], dec["ref"],
	))
}
