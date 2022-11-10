package register

import (
	"encoding/json"
	"net/http"

	"github.com/kxplxn/goteam/server-v2/relay"
)

// Handler is a HTTP handler for the register route.
type Handler struct{}

// NewHandler is the constructor for Handler.
func NewHandler() *Handler {
	return &Handler{}
}

// ServeHTTP responds to requests made to the register route.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// accept only POST
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// decode body into request object
	req := &ReqBody{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		relay.ErrAPIInternal(w, err.Error())
		return
	}

	// create response object to set values into
	res := &ResBody{}

	// validate the request
	if errs := req.Validate(); errs != nil {
		res.Errs = errs
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			relay.ErrAPIInternal(w, err.Error())
		}
		return
	}

	relay.ErrAPIInternal(w, "not implemented")
	return
}
