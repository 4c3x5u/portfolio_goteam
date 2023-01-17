package login

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"server/auth"
	"server/db"
	"server/relay"

	"golang.org/x/crypto/bcrypt"
)

// Handler is the http.Handler for the login route.
type Handler struct {
	userSelector   db.Selector[db.User]
	hashComparer   Comparer
	tokenGenerator auth.TokenGenerator
	dbCloser       db.Closer
}

// NewHandler creates and returns a new Handler.
func NewHandler(
	userSelector db.Selector[db.User],
	hashComparer Comparer,
	tokenGenerator auth.TokenGenerator,
	dbCloser db.Closer,
) Handler {
	return Handler{
		userSelector:   userSelector,
		hashComparer:   hashComparer,
		tokenGenerator: tokenGenerator,
		dbCloser:       dbCloser,
	}
}

// ServeHTTP responds to requests made to the login route.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only accept POST.
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Read and validate request.
	reqBody := ReqBody{}
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
	user, err := h.userSelector.Select(reqBody.Username)
	defer h.dbCloser.Close()
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if err != nil {
		relay.ServerErr(w, err.Error())
		return
	}

	// Compare the password passed in via the request with the hashed password
	// of the user from the database.
	if err = h.hashComparer.Compare(
		user.Password, reqBody.Password,
	); err == bcrypt.ErrMismatchedHashAndPassword {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if err != nil {
		relay.ServerErr(w, err.Error())
		return
	}

	// Generate an authentication cookie for the user and return it within a
	// Set-Cookie header.
	expiry := time.Now().Add(auth.Duration).UTC()
	if authToken, err := h.tokenGenerator.Generate(
		reqBody.Username, expiry,
	); err != nil {
		relay.ServerErr(w, err.Error())
		return
	} else {
		http.SetCookie(w, &http.Cookie{
			Name:    auth.CookieName,
			Value:   authToken,
			Expires: expiry,
		})
		w.WriteHeader(http.StatusOK)
		return
	}
}
