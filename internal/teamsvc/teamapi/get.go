package teamapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/kxplxn/goteam/pkg/cookie"
	"github.com/kxplxn/goteam/pkg/db"
	"github.com/kxplxn/goteam/pkg/db/teamtbl"
	"github.com/kxplxn/goteam/pkg/log"
)

// GetResp defines the body of GET team responses.
type GetResp teamtbl.Team

// GetHandler is an api.MethodHandler that can handle GET requests sent to the
// team route.
type GetHandler struct {
	authDecoder   cookie.Decoder[cookie.Auth]
	teamRetriever db.Retriever[teamtbl.Team]
	teamInserter  db.Inserter[teamtbl.Team]
	teamUpdater   db.Updater[teamtbl.Team]
	inviteEncoder cookie.Encoder[cookie.Invite]
	log           log.Errorer
}

// NewGetHandler creates and returns a new GetHandler.
func NewGetHandler(
	authDecoder cookie.Decoder[cookie.Auth],
	teamRetriever db.Retriever[teamtbl.Team],
	teamInserter db.Inserter[teamtbl.Team],
	teamUpdater db.Updater[teamtbl.Team],
	inviteEncoder cookie.Encoder[cookie.Invite],
	log log.Errorer,
) GetHandler {
	return GetHandler{
		authDecoder:   authDecoder,
		teamRetriever: teamRetriever,
		teamInserter:  teamInserter,
		teamUpdater:   teamUpdater,
		inviteEncoder: inviteEncoder,
		log:           log,
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
	team, err := h.teamRetriever.Retrieve(r.Context(), auth.TeamID)
	var status int
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
		team = teamtbl.NewTeam(
			auth.TeamID,
			[]string{auth.Username},
			[]teamtbl.Board{
				teamtbl.NewBoard(uuid.NewString(), "New Board"),
			},
		)

		// retry a couple of times in the unlitekly event of GUID collision
		for i := 0; i < 3; i++ {
			// insert team into the team table
			if err = h.teamInserter.Insert(
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
		status = http.StatusCreated
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err)
		return
	} else {
		status = http.StatusOK

		if !auth.IsAdmin {
			var isTeamMember bool
			for _, member := range team.Members {
				if member == auth.Username {
					isTeamMember = true
					break
				}
			}
			// if the user is not admin an not a member of the team, add them
			// to the team - this is a synchronisation step and is safe since we
			// validated the JWT and got the username and the team ID from it
			if !isTeamMember {
				team.Members = append(team.Members, auth.Username)
				if err = h.teamUpdater.Update(r.Context(), team); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					h.log.Error(err)
					return
				}
			}

			// return only the boards the user is a member of
			var boards []teamtbl.Board
			for _, b := range team.Boards {
				for _, m := range b.Members {
					if m == auth.Username {
						boards = append(boards, b)
						break
					}
				}
			}
			team.Boards = boards
		}
	}

	// encode invite token if the user is admin
	if auth.IsAdmin {
		ckInv, err := h.inviteEncoder.Encode(cookie.NewInvite(team.ID))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err)
			return
		}
		http.SetCookie(w, &ckInv)
	}

	// encode team
	w.WriteHeader(status)
	if err = json.NewEncoder(w).Encode(GetResp(team)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err)
		return
	}
}
