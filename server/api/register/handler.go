package register

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"server/db"
	"server/relay"

	"github.com/google/uuid"
)

// Handler is the http.Handler for the register route.
type Handler struct {
	validatorReq   Validator
	readerUser     db.Reader[*db.User]
	hasherPwd      Hasher
	creatorUser    db.Creator[*db.User]
	creatorSession db.CreatorTwoStrTime
}

// NewHandler is the constructor for Handler.
func NewHandler(
	validatorReq Validator,
	readerUser db.Reader[*db.User],
	hasherPwd Hasher,
	creatorUser db.Creator[*db.User],
	creatorSession db.CreatorTwoStrTime,
) *Handler {
	return &Handler{
		validatorReq:   validatorReq,
		readerUser:     readerUser,
		hasherPwd:      hasherPwd,
		creatorUser:    creatorUser,
		creatorSession: creatorSession,
	}
}

// ServeHTTP responds to requests made to the register route.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only accept POST.
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Read and validate request.
	reqBody, resBody := &ReqBody{}, &ResBody{}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		relay.ServerErr(w, err.Error())
		return
	}
	if errs := h.validatorReq.Validate(reqBody); errs != nil {
		resBody.Errs = errs
		relay.ClientJSON(w, resBody, http.StatusBadRequest)
		return
	}

	// Check whether the username is taken.
	if _, err := h.readerUser.Read(reqBody.Username); err == nil {
		resBody.Errs = &Errs{Username: []string{errFieldUsernameTaken}}
		relay.ClientJSON(w, resBody, http.StatusBadRequest)
		return
	} else if err != nil && err != sql.ErrNoRows {
		relay.ServerErr(w, err.Error())
		return
	}

	// Hash password and create user.
	if pwdHash, err := h.hasherPwd.Hash(reqBody.Password); err != nil {
		relay.ServerErr(w, err.Error())
		return
	} else if err := h.creatorUser.Create(db.NewUser(reqBody.Username, pwdHash)); err != nil {
		relay.ServerErr(w, err.Error())
		return
	}

	// Create a new session for this user and set session token cookie.
	sessionID := uuid.NewString()
	sessionExpiry := time.Now().Add(1 * time.Hour)
	if err := h.creatorSession.Create(sessionID, reqBody.Username, sessionExpiry); err != nil {
		// User successfuly registered but session creator errored.
		resBody.Errs = &Errs{Session: errSession}
		relay.ClientErr(w, resBody, resBody.Errs.Session, http.StatusUnauthorized)
		return
	} else {
		// Register succes, session creator success, all good...
		http.SetCookie(w, &http.Cookie{
			Name:    "sessionToken",
			Value:   sessionID,
			Expires: sessionExpiry,
		})
		w.WriteHeader(http.StatusOK)
		return
	}
}

// errFieldUsernameTaken is the error message returned from the handler when the
// username given to it is already registered for another user.
const errFieldUsernameTaken = "Username is already taken."

// errFieldUsernameTaken is the error message returned from the handler when
// register is successful but errors occur during session creation.
const errSession = "Register success but session error."
