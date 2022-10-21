package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kxplxn/goteam/server/v2/relay"
)

type HandlerRegister struct {
}

func (h *HandlerRegister) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rly := relay.New(w)

	// accept only POST
	if r.Method != "POST" {
		status := http.StatusMethodNotAllowed
		rly.APIErr(http.StatusText(status), status)
		return
	}

	// read body into map
	dec := make(map[string]string, 3)
	if err := json.NewDecoder(r.Body).Decode(&dec); err != nil {
		status := http.StatusInternalServerError
		rly.APIErr(http.StatusText(status), status)
	}

	// rly decoded body
	rly.APIMsg(fmt.Sprintf(
		"usn: %s\npwd: %s\nref: %s\n",
		dec["usn"], dec["pwd"], dec["ref"],
	))
}
