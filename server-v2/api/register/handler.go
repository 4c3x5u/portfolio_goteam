package register

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/kxplxn/goteam/server-v2/db"
)

// Handler is a HTTP handler for the register endpoint.
type Handler struct{}

// NewHandler is the constructor for Handler handler.
func NewHandler() *Handler {
	return &Handler{}
}

// ServeHTTP responds to requests made to the register endpoint.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// accept only POST
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// decode body into request object
	req := &ReqBody{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
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
			log.Printf("ERROR: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// connect to database
	connStr := os.Getenv(db.ConnStr)
	if connStr == "" {
		log.Print("ERROR: db connection string is empty")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	client, ctx, cancel, err := db.Connect(connStr)
	if err != nil {
		log.Printf("ERROR: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close(client, ctx, cancel)

	// ping database to ensure success
	if err := db.Ping(client, ctx); err != nil {
		log.Printf("ERROR: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// check whether username is unique
	err = client.
		Database("goteamdb").
		Collection("users").
		FindOne(ctx, bson.M{"usn": req.Username}).
		Err()
	if err == nil {
		res.Errs.Username = "Username is already taken."
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Printf("ERROR: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	if err != mongo.ErrNoDocuments {
		log.Printf("ERROR: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
