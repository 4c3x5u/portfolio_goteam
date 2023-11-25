package board

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

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
	teamSelector  dbaccess.Selector[teamTable.Record]
	boardSelector dbaccess.Selector[boardTable.RecursiveRecord]
	log           pkgLog.Errorer
}

// NewGETHandler creates and returns a new GETHandler.
func NewGETHandler(
	userSelector dbaccess.Selector[userTable.Record],
	teamSelector dbaccess.Selector[teamTable.Record],
	boardSelector dbaccess.Selector[boardTable.RecursiveRecord],
	log pkgLog.Errorer,
) GETHandler {
	return GETHandler{
		userSelector:  userSelector,
		teamSelector:  teamSelector,
		boardSelector: boardSelector,
		log:           log,
	}
}

// Handle handles the GET requests sent to the board route.
func (h GETHandler) Handle(
	w http.ResponseWriter, r *http.Request, username string,
) {
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

	_, err = h.teamSelector.Select(strconv.Itoa(user.TeamID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}
}
