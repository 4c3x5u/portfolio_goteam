package login

import (
	"encoding/json"
	"net/http"

	"server/relay"
)

// Handler is the HTTP handler for the login route.
type Handler struct{}

// NewHandler is the constructor for Handler.
func NewHandler() *Handler {
	return &Handler{}
}

// ServeHTTP responds to requests made to the login route. Unlike the register
// handler where we tell the user exactly what's wrong with their credentials,
// we instead just want to return a 400 Bad Request, which the client should
// use to display a boilerplate "Invalid credentials." error.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	reqBody := &ReqBody{}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		relay.ServerErr(w, err.Error())
		return
	}

	if reqBody.Username == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
