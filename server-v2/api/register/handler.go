package register

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

// Handler is a HTTP handler for the register endpoint.
type Handler struct {
	log relay.ErrMsger
}

// NewHandler is the constructor for Handler handler.
func NewHandler(errMsger relay.ErrMsger) *Handler {
	return &Handler{log: errMsger}
}

// ServeHTTP responds to requests made to the register endpoint.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// init logger with writer
	h.log.Init(w)

	// accept only POST
	if r.Method != "POST" {
		h.log.StatusErr(http.StatusMethodNotAllowed)
		return
	}

	// decode body into request object
	req := &ReqBody{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Err(err.Error(), http.StatusBadRequest)
		return
	}

	// create response object
	res := &ResBody{}

	// validate the request
	if isValid, errs := req.IsValid(); !isValid {
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
		res.Errs.Username = "This username is taken."
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
