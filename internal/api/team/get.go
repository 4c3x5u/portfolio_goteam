package team

import (
	"net/http"

	"github.com/kxplxn/goteam/pkg/db"
	teamTable "github.com/kxplxn/goteam/pkg/db/team"
	pkgLog "github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/pkg/token"
)

// GetHandler is an api.MethodHandler that can handle GET requests sent to the
// team route.
type GetHandler struct {
	decodeAuth token.DecodeFunc[token.Auth]
	retriever  db.Retriever[teamTable.Team]
	log        pkgLog.Errorer
}

// NewGetHandler creates and returns a new GetHandler.
func NewGetHandler(
	decodeAuth token.DecodeFunc[token.Auth],
	retriever db.Retriever[teamTable.Team],
	log pkgLog.Errorer,
) GetHandler {
	return GetHandler{
		decodeAuth: decodeAuth,
		retriever:  retriever,
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
	_, err = h.decodeAuth(ckAuth.Value)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}
