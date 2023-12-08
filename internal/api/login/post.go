package login

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/kxplxn/goteam/pkg/db"
	userTable "github.com/kxplxn/goteam/pkg/db/user"
	pkgLog "github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/pkg/token"
)

// POSTReq defines the body of POST login requests.
type POSTReq struct {
	ID       string `json:"username"`
	Password string `json:"password"`
}

// POSTHandler is a http.POSTHandler that can be used to handle login requests.
type POSTHandler struct {
	validator        ReqValidator
	userGetter       db.Getter[userTable.User]
	passwordComparer Comparator
	encodeAuthToken  token.EncodeFunc[token.Auth]
	log              pkgLog.Errorer
}

// NewPOSTHandler creates and returns a new Handler.
func NewPOSTHandler(
	validator ReqValidator,
	userGetter db.Getter[userTable.User],
	hashComparer Comparator,
	encodeAuthToken token.EncodeFunc[token.Auth],
	log pkgLog.Errorer,
) POSTHandler {
	return POSTHandler{
		validator:        validator,
		userGetter:       userGetter,
		passwordComparer: hashComparer,
		encodeAuthToken:  encodeAuthToken,
		log:              log,
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
	user, err := h.userGetter.Get(r.Context(), req.ID)
	if errors.Is(err, db.ErrNoItem) {
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
	exp := time.Now().Add(token.AuthDurationDefault).UTC()
	if authToken, err := h.encodeAuthToken(exp, token.NewAuth(
		user.Username, user.IsAdmin, user.TeamID,
	)); err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		http.SetCookie(w, &http.Cookie{
			Name:     token.AuthName,
			Value:    authToken,
			Expires:  exp,
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
		})
	}
}
