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

// PostReq defines the body of POST login requests.
type PostReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// PostHandler is a http.PostHandler that can be used to handle login requests.
type PostHandler struct {
	validator     ReqValidator
	userRetriever db.Retriever[userTable.User]
	pwdComparator Comparator
	encodeAuth    token.EncodeFunc[token.Auth]
	log           pkgLog.Errorer
}

// NewPostHandler creates and returns a new Handler.
func NewPostHandler(
	validator ReqValidator,
	userRetriever db.Retriever[userTable.User],
	pwdComparator Comparator,
	encodeAuth token.EncodeFunc[token.Auth],
	log pkgLog.Errorer,
) PostHandler {
	return PostHandler{
		validator:     validator,
		userRetriever: userRetriever,
		pwdComparator: pwdComparator,
		encodeAuth:    encodeAuth,
		log:           log,
	}
}

// ServeHTTP responds to requests made to the login route.
func (h PostHandler) Handle(w http.ResponseWriter, r *http.Request, _ string) {
	// Read and validate request body.
	req := PostReq{}
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
	user, err := h.userRetriever.Retrieve(r.Context(), req.Username)
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
	if err = h.pwdComparator.Compare(
		user.Password, req.Password,
	); errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// encode a new auth token
	exp := time.Now().Add(token.DefaultDuration).UTC()
	tkAuth, err := h.encodeAuth(exp, token.NewAuth(
		user.Username, user.IsAdmin, user.TeamID,
	))
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// set auth token in cookie
	http.SetCookie(w, &http.Cookie{
		Name:     token.AuthName,
		Value:    tkAuth,
		Expires:  exp,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})
}
