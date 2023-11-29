package login

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/kxplxn/goteam/server/auth"
	"github.com/kxplxn/goteam/server/dbaccess"
	userTable "github.com/kxplxn/goteam/server/dbaccess/user"
	pkgLog "github.com/kxplxn/goteam/server/log"

	"golang.org/x/crypto/bcrypt"
)

// POSTReq defines the request body for POST login requests.
type POSTReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// POSTHandler is a http.POSTHandler that can be used to handle login requests.
type POSTHandler struct {
	validator          ReqValidator
	userSelector       dbaccess.Selector[userTable.Record]
	passwordComparer   Comparator
	authTokenGenerator auth.TokenGenerator
	log                pkgLog.Errorer
}

// NewPOSTHandler creates and returns a new Handler.
func NewPOSTHandler(
	validator ReqValidator,
	userSelector dbaccess.Selector[userTable.Record],
	hashComparer Comparator,
	authTokenGenerator auth.TokenGenerator,
	log pkgLog.Errorer,
) POSTHandler {
	return POSTHandler{
		validator:          validator,
		userSelector:       userSelector,
		passwordComparer:   hashComparer,
		authTokenGenerator: authTokenGenerator,
		log:                log,
	}
}

// ServeHTTP responds to requests made to the login route.
func (h POSTHandler) Handle(w http.ResponseWriter, r *http.Request, _ string) {
	// Read and validate request body.
	req := POSTReq{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if ok := h.validator.Validate(req); !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Read the user in the database who owns the username that came in the
	// request.
	user, err := h.userSelector.Select(req.Username)
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
		user.Password, req.Password,
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
		req.Username, expiry,
	); err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		http.SetCookie(w, &http.Cookie{
			Name:     auth.CookieName,
			Value:    authToken,
			Expires:  expiry,
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
		})
	}
}
