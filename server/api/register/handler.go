package register

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"server/api"
	"server/auth"
	"server/db"
	pkgLog "server/log"
)

// Handler is a http.Handler that can be used to handle register requests.
type Handler struct {
	validator          ReqValidator
	dbUserSelector     db.Selector[db.User]
	hasher             Hasher
	dbUserInserter     db.Inserter[db.User]
	authTokenGenerator auth.TokenGenerator
	log                pkgLog.Errorer
}

// NewHandler is the constructor for Handler.
func NewHandler(
	validator ReqValidator,
	dbUserSelector db.Selector[db.User],
	hasher Hasher,
	dbUserInserter db.Inserter[db.User],
	authTokenGenerator auth.TokenGenerator,
	log pkgLog.Errorer,
) Handler {
	return Handler{
		validator:          validator,
		dbUserSelector:     dbUserSelector,
		hasher:             hasher,
		dbUserInserter:     dbUserInserter,
		authTokenGenerator: authTokenGenerator,
		log:                log,
	}
}

// ServeHTTP responds to requests made to the register route.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only accept POST.
	if r.Method != http.MethodPost {
		w.Header().Add(api.AllowedMethods(http.MethodPost))
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Read request body.
	reqBody := ReqBody{}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Validate request body, write errors if any occur.
	if validationErrs := h.validator.Validate(reqBody); validationErrs.Any() {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(
			ResBody{ValidationErrs: validationErrs},
		); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	}

	// Check whether the username is taken.
	// TODO: This db call can be removed by adding an "ON CONFLICT (username)
	// DO NOTHING" clause to the query that user inserter uses, and then
	// returning errUsernameTaken if affected rows come back 0. However, not
	// sure if that would increase or decrease the performance as hashing will
	// then occur before exists checks. Test when deployed or when you add
	// integration tests.
	_, err := h.dbUserSelector.Select(reqBody.Username)
	if err == nil {
		w.WriteHeader(http.StatusBadRequest)
		if errEncode := json.NewEncoder(w).Encode(
			ResBody{ValidationErrs: ValidationErrs{
				Username: []string{"Username is already taken."},
			}},
		); errEncode != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(errEncode.Error())
		}
		return
	} else if err != sql.ErrNoRows {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Hash password and create user.
	if pwdHash, err := h.hasher.Hash(reqBody.Password); err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if err = h.dbUserInserter.Insert(
		db.NewUser(reqBody.Username, pwdHash),
	); err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Generate an authentication token for the user that is valid for an hour
	// and return it within a Set-Cookie header.
	expiry := time.Now().Add(auth.Duration).UTC()
	if authToken, err := h.authTokenGenerator.Generate(
		reqBody.Username, expiry,
	); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(
			ResBody{
				Msg: "You have been registered successfully but something " +
					"went wrong. Please log in using the credentials you " +
					"registered with.",
			},
		); err != nil {
			h.log.Error(err.Error())
		}
	} else {
		http.SetCookie(w, &http.Cookie{
			Name:    auth.CookieName,
			Value:   authToken,
			Expires: expiry,
		})
		w.WriteHeader(http.StatusOK)
	}
}
