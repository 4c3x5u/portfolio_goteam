package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type RegisterHandler struct {
}

func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	relay := NewLogger(w)
	relay.Msg("register endpoint is hit")

	switch r.Method {
	case "POST":
		h.POST(relay, r)
	default:
		relay.Err("Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *RegisterHandler) POST(relay *Logger, request *http.Request) {
	body := make(map[string]string)
	if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
		statusCode := http.StatusInternalServerError
		relay.Err(http.StatusText(statusCode), statusCode)
	}

	relay.Msg(fmt.Sprintf("usn: %s\npwd: %s", body["usn"], body["pwd"]))
}
