package login

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"server/auth"
	"server/db"
	"server/relay"
)

// Handler is the HTTP handler for the login route.
type Handler struct {
	readerUser     db.Reader[*db.User]
	comparerHash   Comparer
	generatorToken auth.Generator
}

// NewHandler is the constructor for Handler.
func NewHandler(
	readerUser db.Reader[*db.User],
	comparerHash Comparer,
	generatorToken auth.Generator,
) *Handler {
	return &Handler{
		readerUser:     readerUser,
		comparerHash:   comparerHash,
		generatorToken: generatorToken,
	}
}

// ServeHTTP responds to requests made to the login route. Unlike the register
// handler where we tell the user exactly what's wrong with their credentials,
// we instead just want to return a 400 Bad Request, which the client should
// use to display a boilerplate "Invalid credentials." error.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only accept POST.
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Read and validate request.
	reqBody := &ReqBody{}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		relay.ServerErr(w, err.Error())
		return
	}
	if reqBody.Username == "" || reqBody.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Read the user in the database who owns the username that came in the
	// request.
	user, err := h.readerUser.Read(reqBody.Username)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if err != nil {
		relay.ServerErr(w, err.Error())
		return
	}

	// Compare the password passed in via the request with the hashed password
	// of the user from the database.
	if isMatch, err := h.comparerHash.Compare(user.Password, reqBody.Password); err != nil {
		relay.ServerErr(w, err.Error())
		return
	} else if !isMatch {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Generate a JWT for the user and return it in a Set-Cookie header
	expiry := time.Now().Add(1 * time.Hour)
	if tokenStr, err := h.generatorToken.Generate(reqBody.Username, expiry); err != nil {
		relay.ServerErr(w, err.Error())
		return
	} else {
		http.SetCookie(w, &http.Cookie{
			Name:    "authToken",
			Value:   tokenStr,
			Expires: expiry.UTC(),
		})
		w.WriteHeader(http.StatusOK)
		return
	}
}
