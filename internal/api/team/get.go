package team

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/kxplxn/goteam/pkg/db"
	teamTable "github.com/kxplxn/goteam/pkg/db/team"
	pkgLog "github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/pkg/token"
)

// GetResp defines the body of GET team responses.
type GetResp teamTable.Team

// GetHandler is an api.MethodHandler that can handle GET requests sent to the
// team route.
type GetHandler struct {
	decodeAuth token.DecodeFunc[token.Auth]
	retriever  db.Retriever[teamTable.Team]
	inserter   db.Inserter[teamTable.Team]
	log        pkgLog.Errorer
}

// NewGetHandler creates and returns a new GetHandler.
func NewGetHandler(
	decodeAuth token.DecodeFunc[token.Auth],
	retriever db.Retriever[teamTable.Team],
	inserter db.Inserter[teamTable.Team],
	log pkgLog.Errorer,
) GetHandler {
	return GetHandler{
		decodeAuth: decodeAuth,
		retriever:  retriever,
		inserter:   inserter,
		log:        log,
	}
}

// Handle handles GET requests sent to the team route.
func (h GetHandler) Handle(w http.ResponseWriter, r *http.Request, _ string) {
	// get auth token
	ckAuth, err := r.Cookie(token.AuthName)
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// decode auth token
	auth, err := h.decodeAuth(ckAuth.Value)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// retrieve team
	team, err := h.retriever.Retrieve(r.Context(), auth.TeamID)
	if errors.Is(err, db.ErrNoItem) {
		// if team was not found and since we trust the JWT, this is our sign
		// from the register endpoint that this is a new user and we should
		// create a new team for them

		// register endpoint must have set the isAdmin to true
		// this check might be redundant but it's here just in case
		if !auth.IsAdmin {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// create team
		team = teamTable.NewTeam(
			auth.TeamID,
			[]string{auth.Username},
			[]teamTable.Board{
				teamTable.NewBoard(uuid.NewString(), "New Board"),
			},
		)

		// insert team into the team table
		if err = h.inserter.Insert(r.Context(), team); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
			return
		}

		// write 201 to indicate creation of the new team
		w.WriteHeader(http.StatusCreated)
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}

	// encode team
	if err = json.NewEncoder(w).Encode(GetResp(team)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}
}
