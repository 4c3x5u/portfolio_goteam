package register

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"server/auth"
	"server/db"
	"server/relay"
)

// Handler is the http.Handler for the register route.
type Handler struct {
	validator      Validator
	userSelector   db.Selector[db.User]
	hasher         Hasher
	userInserter   db.Inserter[db.User]
	tokenGenerator auth.TokenGenerator
	dbCloser       db.Closer
}

// NewHandler is the constructor for Handler.
func NewHandler(
	validator Validator,
	userSelector db.Selector[db.User],
	hasher Hasher,
	userInserter db.Inserter[db.User],
	tokenGenerator auth.TokenGenerator,
	dbCloser db.Closer,
) Handler {
	return Handler{
		validator:      validator,
		userSelector:   userSelector,
		hasher:         hasher,
		userInserter:   userInserter,
		tokenGenerator: tokenGenerator,
		dbCloser:       dbCloser,
	}
}

// ServeHTTP responds to requests made to the register route.
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
	if errs := h.validator.Validate(reqBody); errs.Any() {
		relay.ClientJSON(w, http.StatusBadRequest, ResBody{ValidationErrs: errs})
		return
	}

	// Check whether the username is taken. This db call can be removed by
	// adding an "ON CONFLICT (username) DO NOTHING" clause to the query that
	// user inserter uses, and then returning errUsernameTaken if affected
	// rows come back 0. However, not sure if that would increase or decrease
	// the performance as hashing will then occur before exists checks.
	// TODO: Test when deployed.
	_, err := h.userSelector.Select(reqBody.Username)
	defer h.dbCloser.Close()
	if err == nil {
		relay.ClientJSON(w, http.StatusBadRequest, ResBody{
			ValidationErrs: ValidationErrs{Username: []string{errUsernameTaken}},
		})
		return
	} else if err != sql.ErrNoRows {
		relay.ServerErr(w, err.Error())
		return
	}

	// Hash password and create user.
	if pwdHash, err := h.hasher.Hash(reqBody.Password); err != nil {
		relay.ServerErr(w, err.Error())
		return
	} else if err = h.userInserter.Insert(
		db.NewUser(reqBody.Username, pwdHash),
	); err != nil {
		relay.ServerErr(w, err.Error())
		return
	}

	// Generate an authentication token for the user that is valid for an hour
	// and return it within a Set-Cookie header.
	expiry := time.Now().Add(auth.Duration).UTC()
	if authToken, err := h.tokenGenerator.Generate(
		reqBody.Username, expiry,
	); err != nil {
		relay.ClientErr(w, http.StatusUnauthorized, ResBody{
			ValidationErrs: ValidationErrs{Auth: errAuth},
		})
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

// errUsernameTaken is the error message returned from the handler when the
// username given to it is already registered for another user.
const errUsernameTaken = "Username is already taken."

// errAuth is the error message returned from handlers when the token generator
// throws an error
const errAuth = "You have been registered successfully but something went wrong. " +
	"Please log in using the credentials you registered with."
