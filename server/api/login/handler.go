package login

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"server/db"
	"server/relay"
)

// Handler is the HTTP handler for the login route.
type Handler struct {
	readerUser    db.Reader[*db.User]
	comparerHash  Comparer
	readerSession db.Reader[*db.Session]
}

// NewHandler is the constructor for Handler.
func NewHandler(
	readerUser db.Reader[*db.User],
	comparerHash Comparer,
	readerSession db.Reader[*db.Session],
) *Handler {
	return &Handler{
		readerUser:    readerUser,
		comparerHash:  comparerHash,
		readerSession: readerSession,
	}
}

// ServeHTTP responds to requests made to the login route. Unlike the register
// handler where we tell the user exactly what's wrong with their credentials,
// we instead just want to return a 400 Bad Request, which the client should
// use to display a boilerplate "Invalid credentials." error.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	reqBody := &ReqBody{}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		relay.ServerErr(w, err.Error())
		return
	}

	if reqBody.Username == "" || reqBody.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.readerUser.Read(reqBody.Username)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if err != nil {
		relay.ServerErr(w, err.Error())
		return
	}

	isMatch, err := h.comparerHash.Compare(user.Password, reqBody.Password)
	if err != nil {
		relay.ServerErr(w, err.Error())
		return
	}
	if !isMatch {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = h.readerSession.Read(reqBody.Username)
	if err != nil && err != sql.ErrNoRows {
		relay.ServerErr(w, err.Error())
		return
	}
}
