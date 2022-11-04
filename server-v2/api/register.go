package api

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"os"

	"github.com/kxplxn/goteam/server-v2/db"
	"github.com/kxplxn/goteam/server-v2/relay"
)

// ReqRegister is the request contract for the register endpoint.
type ReqRegister struct {
	Usn string `json:"username"`
	Pwd string `json:"password"`
	Ref string `json:"referrer"`
}

// ResRegister defines the resposne type for the register endpoint.
type ResRegister struct {
	Errs *ErrsRegister `json:"errors"`
}

// ErrsRegister defines the structure of error object that can be encoded in the
// register endpoint in the case of an error.
type ErrsRegister struct {
	Usn []string `json:"username"`
}

// Validate checks whether all the error fields are empty on the ErrsRegister
// object — returns true if so and false otherwise. If more error fields are
// added to ErrsRegister, they should also
func (e *ErrsRegister) Validate() bool {
	return len(e.Usn) == 0
}

// HandlerRegister is a HTTP handler for the register endpoint.
type HandlerRegister struct {
	log relay.ErrMsger
}

// NewHandlerRegister is the constructor for HandlerRegister handler.
func NewHandlerRegister(errMsger relay.ErrMsger) *HandlerRegister {
	return &HandlerRegister{log: errMsger}
}

// ServeHTTP responds to requests made to the to the register endpoint.
func (h *HandlerRegister) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// init logger with writer
	h.log.Init(w)

	// accept only POST
	if r.Method != "POST" {
		h.log.StatusErr(http.StatusMethodNotAllowed)
		return
	}

	// decode body into request type
	req := &ReqRegister{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Err(err.Error(), http.StatusInternalServerError)
		return
	}

	// set up database
	connStr := os.Getenv(db.ConnStr)
	if connStr == "" {
		h.log.Err("db connection string empty", http.StatusInternalServerError)
		return
	}
	client, ctx, cancel, err := db.Connect(connStr)
	if err != nil {
		h.log.Err(err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close(client, ctx, cancel)

	// ping database to ensure success
	if err := db.Ping(client, ctx); err != nil {
		h.log.Err(err.Error(), http.StatusInternalServerError)
		return
	}

	res := &ResRegister{}

	// check whether username is unique
	err = client.
		Database("goteamdb").
		Collection("users").
		FindOne(ctx, bson.M{"usn": req.Usn}).
		Err()
	if err == nil {
		h.log.Err("username taken", http.StatusBadRequest)
		res.Errs.Usn = append(res.Errs.Usn, "This username is taken.")
		return
	}
	if err != mongo.ErrNoDocuments {
		h.log.Err(err.Error(), http.StatusInternalServerError)
		return
	}

	// return errors if exist
	if res.Errs.Validate() != true {
		h.log.ResErr(res, "errors on register", http.StatusBadRequest)
		return
	}

	// relay request fields – for testing purposes
	h.log.Msg(fmt.Sprintf(
		"usn: %s\npwd: %s\nref: %s\n",
		req.Usn, req.Pwd, req.Ref,
	))
}
