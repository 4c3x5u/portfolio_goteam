package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/kxplxn/goteam/server-v2/db"
	"github.com/kxplxn/goteam/server-v2/relay"
)

// ReqRegister is the request contract for the register endpoint.
type ReqRegister struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Referrer string `json:"referrer"`
}

// Validate performs input validation checks on the register request and returns
// an error if any fails.
func (r *ReqRegister) Validate() (isErr bool, errs *ErrsRegister) {
	errs = &ErrsRegister{}

	// username too short
	if len(r.Username) < 5 {
		errs.Username = append(errs.Username, "Username cannot be shorter than 5 characters.")
		return true, errs
	}

	return false, nil
}

// ErrsRegister defines the structure of error object that can be encoded in the
// register endpoint in the case of an error.
type ErrsRegister struct {
	Username []string `json:"username"`
}

// ResRegister defines the resposne type for the register endpoint.
type ResRegister struct {
	Errs *ErrsRegister `json:"errors"`
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

	// decode body into request object
	req := &ReqRegister{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Err(err.Error(), http.StatusBadRequest)
		return
	}

	// create response object
	res := &ResRegister{}

	// validate the request
	if isErr, errs := req.Validate(); isErr {
		res.Errs = errs
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			h.log.Err(err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// connect to database
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

	// check whether username is unique
	err = client.
		Database("goteamdb").
		Collection("users").
		FindOne(ctx, bson.M{"usn": req.Username}).
		Err()
	if err == nil {
		h.log.Err("username taken", http.StatusBadRequest)
		res.Errs.Username = append(res.Errs.Username, "This username is taken.")
		return
	}
	if err != mongo.ErrNoDocuments {
		h.log.Err(err.Error(), http.StatusInternalServerError)
		return
	}

	// todo: remove
	// relay request fields – for testing purposes
	h.log.Msg(fmt.Sprintf(
		"usn: %s\npwd: %s\nref: %s\n",
		req.Username, req.Password, req.Referrer,
	))
}
