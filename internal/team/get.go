package team

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/kxplxn/goteam/pkg/cookie"
	"github.com/kxplxn/goteam/pkg/db"
	"github.com/kxplxn/goteam/pkg/db/teamtable"
	"github.com/kxplxn/goteam/pkg/log"
)

// GetResp defines the body of GET team responses.
type GetResp teamtable.Team

// GetHandler is an api.MethodHandler that can handle GET requests sent to the
// team route.
type GetHandler struct {
	authDecoder cookie.Decoder[cookie.Auth]
	retriever   db.Retriever[teamtable.Team]
	inserter    db.Inserter[teamtable.Team]
	log         log.Errorer
}

// NewGetHandler creates and returns a new GetHandler.
func NewGetHandler(
	decodeAuth cookie.Decoder[cookie.Auth],
	retriever db.Retriever[teamtable.Team],
	inserter db.Inserter[teamtable.Team],
	log log.Errorer,
) GetHandler {
	return GetHandler{
		authDecoder: decodeAuth,
		retriever:   retriever,
		inserter:    inserter,
		log:         log,
	}
}

// Handle handles GET requests sent to the team route.
func (h GetHandler) Handle(w http.ResponseWriter, r *http.Request, _ string) {
	// get auth token
	ckAuth, err := r.Cookie(cookie.AuthName)
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// decode auth token
	auth, err := h.authDecoder.Decode(*ckAuth)
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
		team = teamtable.NewTeam(
			auth.TeamID,
			[]string{auth.Username},
			[]teamtable.Board{
				teamtable.NewBoard(uuid.NewString(), "New Board"),
			},
		)

		// retry a couple of times in the unlitekly event of GUID collision
		for i := 0; i < 3; i++ {
			// insert team into the team table
			if err = h.inserter.Insert(
				r.Context(), team,
			); errors.Is(err, db.ErrDupKey) {
				team.Boards[0].ID = uuid.NewString()
			} else if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				h.log.Error(err)
				return
			} else {
				break
			}
		}

		// write 201 to indicate creation of the new team
		w.WriteHeader(http.StatusCreated)
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err)
		return
	}

	// encode team
	if err = json.NewEncoder(w).Encode(GetResp(team)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err)
		return
	}
}
