package board

import (
	"database/sql"
	"encoding/json"
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

	// Select recursive board record from the database.
	activeBoard, err := h.boardSelector.Select(boardID)
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}

	// Validate that the user is admin or member of the team that the board
	// belongs to.
	if activeBoard.TeamID != user.TeamID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Select the team from the database that the user is the member/admin of.
	team, err := h.teamSelector.Select(strconv.Itoa(user.TeamID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}

	// Build response from data retrieved from the database.
	resp := GETResp{
		Username: username,
		Team: Team{
			ID:         team.ID,
			InviteCode: team.InviteCode,
		},
		// TODO: get all team members
		TeamMembers: []TeamMember{},
		// TODO: get all boards
		Boards: []Board{},
		ActiveBoard: ActiveBoard{
			ID:      activeBoard.ID,
			Name:    activeBoard.Name,
			Columns: make([]Column, len(activeBoard.Columns)),
		},
	}
	for i, col := range activeBoard.Columns {
		resp.ActiveBoard.Columns[i].ID = col.ID
		resp.ActiveBoard.Columns[i].Order = col.Order
		resp.ActiveBoard.Columns[i].Tasks = make([]Task, len(col.Tasks))

		for j, task := range col.Tasks {
			resp.ActiveBoard.Columns[i].Tasks[j].ID = task.ID
			resp.ActiveBoard.Columns[i].Tasks[j].Title = task.Title
			resp.ActiveBoard.Columns[i].Tasks[j].Description = task.Description
			resp.ActiveBoard.Columns[i].Tasks[j].Order = task.Order
			resp.ActiveBoard.Columns[i].Tasks[j].Subtasks = make(
				[]Subtask, len(task.Subtasks),
			)

			for k, subtask := range task.Subtasks {
				resp.ActiveBoard.Columns[i].Tasks[j].Subtasks[k].
					ID = subtask.ID
				resp.ActiveBoard.Columns[i].Tasks[j].Subtasks[k].
					Title = subtask.Title
				resp.ActiveBoard.Columns[i].Tasks[j].Subtasks[k].
					Order = subtask.Order
				resp.ActiveBoard.Columns[i].Tasks[j].Subtasks[k].
					IsDone = subtask.IsDone
			}
		}
	}

	// Return response body, or 500 if JSON encoding fails.
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}
}

// GETResp defines the response body for GET board requests.
type GETResp struct {
	Username    string       `json:"user"`
	Team        Team         `json:"team"`
	TeamMembers []TeamMember `json:"members"`
	Boards      []Board      `json:"boards"`
	ActiveBoard ActiveBoard  `json:"activeBoard"`
}

// Team defines the team data returned in GETResp.
type Team struct {
	ID         int    `json:"id"`
	InviteCode string `json:"inviteCode"`
}

// TeamMember defines an item in the team members data returned in GETResp.
type TeamMember struct {
	Username string `json:"username"`
	IsAdmin  string `json:"isAdmin"`
	// TODO: figure out what isActive from client is and decide what to do
}

// Board defines an item in the boards data returned in GETResp.
type Board struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ActiveBoard defines the active board data returned in GETResp.
type ActiveBoard struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Columns []Column `json:"columns"`
}

// Column defines an item in the column data returned in ActiveBoard.
type Column struct {
	ID    int    `json:"id"`
	Order int    `json:"order"`
	Tasks []Task `json:"tasks"`
}

// Task defines an item in the task data returned in Column.
type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Order       int       `json:"order"`
	Subtasks    []Subtask `json:"subtasks"`
}

// Subtask defines an item in the subtask data returned in Task.
type Subtask struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Order  int    `json:"order"`
	IsDone bool   `json:"isDone"`
}
