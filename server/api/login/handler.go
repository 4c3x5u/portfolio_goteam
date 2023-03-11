package login

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"server/api"
	"server/auth"
	"server/db"
	"server/log"

	"golang.org/x/crypto/bcrypt"
)

// Handler is a http.Handler that can be used to handle login requests.
type Handler struct {
	validator          ReqValidator
	dbUserSelector     db.Selector[db.User]
	passwordComparer   Comparer
	authTokenGenerator auth.TokenGenerator
	logger             log.Logger
}

// NewHandler creates and returns a new Handler.
func NewHandler(
	validator ReqValidator,
	userSelector db.Selector[db.User],
	hashComparer Comparer,
	authTokenGenerator auth.TokenGenerator,
	logger log.Logger,
) Handler {
	return Handler{
		validator:          validator,
		dbUserSelector:     userSelector,
		passwordComparer:   hashComparer,
		authTokenGenerator: authTokenGenerator,
		logger:             logger,
	}
}

// ServeHTTP responds to requests made to the login route.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only accept POST.
	if r.Method != http.MethodPost {
		w.Header().Add(api.AllowedMethods(http.MethodPost))
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Read and validate request body.
	reqBody := ReqBody{}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		h.logger.Log(log.LevelError, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if ok := h.validator.Validate(reqBody); !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Read the user in the database who owns the username that came in the
	// request.
	user, err := h.dbUserSelector.Select(reqBody.Username)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if err != nil {
		h.logger.Log(log.LevelError, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Compare the password passed in via the request with the hashed password
	// of the user from the database.
	if err = h.passwordComparer.Compare(
		user.Password, reqBody.Password,
	); err == bcrypt.ErrMismatchedHashAndPassword {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if err != nil {
		h.logger.Log(log.LevelError, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Generate an authentication cookie for the user and return it within a
	// Set-Cookie header.
	expiry := time.Now().Add(auth.Duration).UTC()
	if authToken, err := h.authTokenGenerator.Generate(
		reqBody.Username, expiry,
	); err != nil {
		h.logger.Log(log.LevelError, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		http.SetCookie(w, &http.Cookie{
			Name:    auth.CookieName,
			Value:   authToken,
			Expires: expiry,
		})
		w.WriteHeader(http.StatusOK)
	}
}
