package board

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/kxplxn/goteam/server/api"
	"github.com/kxplxn/goteam/server/dbaccess"
	boardTable "github.com/kxplxn/goteam/server/dbaccess/board"
	teamTable "github.com/kxplxn/goteam/server/dbaccess/team"
	userTable "github.com/kxplxn/goteam/server/dbaccess/user"
	pkgLog "github.com/kxplxn/goteam/server/log"
)

// GETHandler is an api.MethodHandler that can be used to handle GET board
// requests.
type GETHandler struct {
	userSelector  dbaccess.Selector[userTable.Record]
	idValidator   api.StringValidator
	boardSelector dbaccess.Selector[boardTable.RecursiveRecord]
	teamSelector  dbaccess.Selector[teamTable.Record]
	log           pkgLog.Errorer
}

// NewGETHandler creates and returns a new GETHandler.
func NewGETHandler(
	userSelector dbaccess.Selector[userTable.Record],
	idValidator api.StringValidator,
	boardSelector dbaccess.Selector[boardTable.RecursiveRecord],
	teamSelector dbaccess.Selector[teamTable.Record],
	log pkgLog.Errorer,
) GETHandler {
	return GETHandler{
		userSelector:  userSelector,
		idValidator:   idValidator,
		boardSelector: boardSelector,
		teamSelector:  teamSelector,
		log:           log,
	}
}

// Handle handles the GET requests sent to the board route.
func (h GETHandler) Handle(
	w http.ResponseWriter, r *http.Request, username string,
) {
	// Select the user from the database to access their TeamID.
	user, err := h.userSelector.Select(username)
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusUnauthorized)
		h.log.Error(err.Error())
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}

	// Validate board ID.
	boardID := r.URL.Query().Get("id")
	if err = h.idValidator.Validate(boardID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = h.boardSelector.Select(boardID)
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}

	// Select the team from the database that the user is the member/admin of.
	_, err = h.teamSelector.Select(strconv.Itoa(user.TeamID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}
}
