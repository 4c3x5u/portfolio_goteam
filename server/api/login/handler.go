package login

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/kxplxn/goteam/server/api"
	"github.com/kxplxn/goteam/server/auth"
	"github.com/kxplxn/goteam/server/dbaccess"
	userTable "github.com/kxplxn/goteam/server/dbaccess/user"
	pkgLog "github.com/kxplxn/goteam/server/log"

	"golang.org/x/crypto/bcrypt"
)

// Handler is a http.Handler that can be used to handle login requests.
type Handler struct {
	validator          ReqValidator
	userSelector       dbaccess.Selector[userTable.Record]
	passwordComparer   Comparator
	authTokenGenerator auth.TokenGenerator
	log                pkgLog.Errorer
}

// NewHandler creates and returns a new Handler.
func NewHandler(
	validator ReqValidator,
	userSelector dbaccess.Selector[userTable.Record],
	hashComparer Comparator,
	authTokenGenerator auth.TokenGenerator,
	log pkgLog.Errorer,
) Handler {
	return Handler{
		validator:          validator,
		userSelector:       userSelector,
		passwordComparer:   hashComparer,
		authTokenGenerator: authTokenGenerator,
		log:                log,
	}
}

// ServeHTTP responds to requests made to the login route.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only accept POST.
	if r.Method != http.MethodPost {
		w.Header().Add(api.AllowedMethods([]string{http.MethodPost}))
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Read and validate request body.
	reqBody := ReqBody{}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if ok := h.validator.Validate(reqBody); !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Read the user in the database who owns the username that came in the
	// request.
	user, err := h.userSelector.Select(reqBody.Username)
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Compare the password passed in via the request with the hashed password
	// of the user from the database.
	if err = h.passwordComparer.Compare(
		user.Password, reqBody.Password,
	); errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Generate an authentication cookie for the user and return it within a
	// Set-Cookie header.
	expiry := time.Now().Add(auth.Duration).UTC()
	if authToken, err := h.authTokenGenerator.Generate(
		reqBody.Username, expiry,
	); err != nil {
		h.log.Error(err.Error())
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
