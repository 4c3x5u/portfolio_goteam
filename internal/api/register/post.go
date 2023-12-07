package register

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kxplxn/goteam/pkg/auth"
	"github.com/kxplxn/goteam/pkg/db"
	userTable "github.com/kxplxn/goteam/pkg/db/user"
	pkgLog "github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/pkg/token"
)

// PostReq defines the body of POST register requests.
type PostReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	TeamID   string `json:"teamID"`
}

// PostResp defines the body of POST register responses.
type PostResp struct {
	Err            string         `json:"error,omitempty"`
	ValidationErrs ValidationErrs `json:"validationErrors,omitempty"`
}

// PostHandler is a api.MethodHandler that can be used to handle POST register
// requests.
type PostHandler struct {
	reqValidator      ReqValidator
	hasher            Hasher
	decodeInviteToken token.DecodeFunc[token.Invite]
	userPutter        db.Putter[userTable.User]
	encodeAuthToken   token.EncodeFunc[token.Auth]
	log               pkgLog.Errorer
}

// NewPostHandler creates and returns a new HandlerPost.
func NewPostHandler(
	userValidator ReqValidator,
	decodeInviteToken token.DecodeFunc[token.Invite],
	hasher Hasher,
	userPutter db.Putter[userTable.User],
	encodeAuthToken token.EncodeFunc[token.Auth],
	log pkgLog.Errorer,
) PostHandler {
	return PostHandler{
		reqValidator:      userValidator,
		hasher:            hasher,
		decodeInviteToken: decodeInviteToken,
		userPutter:        userPutter,
		encodeAuthToken:   encodeAuthToken,
		log:               log,
	}
}

// ServeHTTP responds to requests made to the register route.
func (h PostHandler) Handle(w http.ResponseWriter, r *http.Request, _ string) {
	// decode request
	req := PostReq{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// validate request
	errsValidate := h.reqValidator.Validate(req)
	if errsValidate.Any() {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(
			PostResp{ValidationErrs: errsValidate},
		); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	}

	// create user
	user := userTable.User{
		Username: req.Username,
		Password: []byte{},
		IsAdmin:  false,
		TeamID:   "",
	}

	// set user's TeamID and IsAdmin based on invite token.
	ck, err := r.Cookie(token.NameInvite)
	if err == http.ErrNoCookie {
		user.TeamID = uuid.NewString()
		user.IsAdmin = true
	} else if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		invite, err := h.decodeInviteToken(ck.Value)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(
				PostResp{Err: "Invalid invite token."},
			); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				h.log.Error(err.Error())
			}
			return
		}
		user.TeamID = invite.TeamID
		user.IsAdmin = false
	}

	// hash password
	pwdHash, err := h.hasher.Hash(req.Password)
	if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user.Password = pwdHash

	// put user into the user table
	if err = h.userPutter.Put(r.Context(), user); err == db.ErrDupKey {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(
			PostResp{ValidationErrs: ValidationErrs{
				Username: []string{"Username is already taken."},
			}},
		); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	} else if err != nil {
		h.log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// generate an auth token
	exp := time.Now().Add(auth.Duration).UTC()
	tkAuth, err := h.encodeAuthToken(exp, token.NewAuth(
		user.Username, user.IsAdmin, user.TeamID,
	))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(
			PostResp{
				Err: "You have been registered successfully but something " +
					"went wrong. Please log in using the credentials you " +
					"registered with.",
			},
		); err != nil {
			h.log.Error(err.Error())
		}
		return
	}

	// set auth cookie and respond OK
	http.SetCookie(w, &http.Cookie{
		Name:     token.NameAuth,
		Value:    tkAuth,
		Expires:  exp,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})
	w.WriteHeader(http.StatusOK)
}
