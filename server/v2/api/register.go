package api

import (
	"encoding/json"
	"fmt"
	"github.com/kxplxn/goteam/server/v2/log"
	"net/http"
)

func ServeRegister(w http.ResponseWriter, r *http.Request) {
	relay := log.NewAPILogger(w)
	relay.Msg("register endpoint is hit")

	switch r.Method {
	case "POST":
		post(relay, r)
	default:
		relay.Err("Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func post(relay *log.APILogger, request *http.Request) {
	body := make(map[string]string)
	if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
		statusCode := http.StatusInternalServerError
		relay.Err(http.StatusText(statusCode), statusCode)
	}

	relay.Msg(fmt.Sprintf("usn: %s\npwd: %s", body["usn"], body["pwd"]))
}
