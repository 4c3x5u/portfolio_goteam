package register

import (
	"encoding/json"
	"net/http"

	"github.com/kxplxn/goteam/server-v2/relay"
)

// Handler is a HTTP handler for the register route.
type Handler struct {
	creatorUser CreatorUser
}

// NewHandler is the constructor for Handler.
func NewHandler(creatorUser CreatorUser) *Handler {
	return &Handler{creatorUser: creatorUser}
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
		relay.ServerErr(w, err.Error())
		return
	}

	// create response object to set values into
	res := &ResBody{}

	// validate the request
	if errs := req.Validate(); errs != nil {
		res.ErrsValidation = errs
		relay.ClientJSON(w, res, http.StatusBadRequest)
		return
	}

	errsValidation, err := h.creatorUser.CreateUser(string(req.Username), string(req.Password))
	if err != nil {
		relay.ServerErr(w, err.Error())
		return
	}
	if errsValidation != nil {
		res.ErrsValidation = errsValidation
		relay.ClientJSON(w, res, http.StatusBadRequest)
		return
	}

	relay.ServerErr(w, "not implemented")
	return
}
