package api

import (
	"encoding/json"
	"fmt"
	"github.com/kxplxn/goteam/server/v2/db"
	"github.com/kxplxn/goteam/server/v2/relay"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"os"
)

// ReqRegister is the request contract for the register endpoint.
type ReqRegister struct {
	Usn string `json:"usn"` // username
	Pwd string `json:"pwd"` // password
	Ref string `json:"ref"` // referrer
}

// HandlerRegister is a HTTP handler for the register endpoint.
type HandlerRegister struct {
	log relay.APIErrMsger
}

// NewHandlerRegister is the constructor for HandlerRegister handler.
func NewHandlerRegister(errMsger relay.APIErrMsger) *HandlerRegister {
	return &HandlerRegister{log: errMsger}
}

// todo: are "return"s in error blocks necessary/verbose?

// ServeHTTP responds to requests made to the to the register endpoint.
func (h *HandlerRegister) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// accept only POST
	if r.Method != "POST" {
		h.log.ErrStatus(w, http.StatusMethodNotAllowed)
		return
	}

	// decode body into request type
	req := &ReqRegister{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Err(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// set up database
	connStr := os.Getenv(db.ConnStr)
	if connStr == "" {
		h.log.Err(w, "db connection string empty", http.StatusInternalServerError)
		return
	}
	client, ctx, cancel, err := db.Connect(connStr)
	if err != nil {
		h.log.Err(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close(client, ctx, cancel)

	// check whether username is unique
	err = client.
		Database("goteamdb").
		Collection("users").
		FindOne(ctx, bson.M{"usn": req.Usn}).
		Err()
	if err == nil {
		// todo: better validation errors
		h.log.Err(w, "username taken", http.StatusBadRequest)
		return
	}
	if err != mongo.ErrNoDocuments {
		h.log.Err(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// relay request fields
	h.log.Msg(w, fmt.Sprintf(
		"usn: %s\npwd: %s\nref: %s\n",
		req.Usn, req.Pwd, req.Ref,
	))
}
