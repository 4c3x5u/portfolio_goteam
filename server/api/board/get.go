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

// GETResp defines the response body for GET board requests.
type GETResp struct {
	User        User         `json:"user"`
	Team        Team         `json:"team"`
	TeamMembers []TeamMember `json:"members"`
	Boards      []Board      `json:"boards"`
	ActiveBoard ActiveBoard  `json:"activeBoard"`
}

// User defines the user data return in GetResp.
type User struct {
	Username string `json:"username"`
	IsAdmin  bool   `json:"isAdmin"`
}

// Team defines the team data returned in GETResp.
type Team struct {
	ID         int    `json:"id"`
	InviteCode string `json:"inviteCode"`
}

// TeamMember defines an item in the team members data returned in GETResp.
type TeamMember struct {
	Username string `json:"username"`
	IsAdmin  bool   `json:"isAdmin"`
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

// GETHandler is an api.MethodHandler that can be used to handle GET board
// requests.
type GETHandler struct {
	userSelector           dbaccess.Selector[userTable.Record]
	boardInserter          dbaccess.Inserter[boardTable.InRecord]
	idValidator            api.StringValidator
	boardSelectorRecursive dbaccess.Selector[boardTable.RecursiveRecord]
	teamSelector           dbaccess.Selector[teamTable.Record]
	userSelectorByTeamID   dbaccess.Selector[[]userTable.Record]
	boardSelectorByTeamID  dbaccess.Selector[[]boardTable.Record]
	log                    pkgLog.Errorer
}

// NewGETHandler creates and returns a new GETHandler.
func NewGETHandler(
	userSelector dbaccess.Selector[userTable.Record],
	boardInserter dbaccess.Inserter[boardTable.InRecord],
	idValidator api.StringValidator,
	boardSelectorRecursive dbaccess.Selector[boardTable.RecursiveRecord],
	teamSelector dbaccess.Selector[teamTable.Record],
	userSelectorByTeamID dbaccess.Selector[[]userTable.Record],
	boardSelectorByTeamID dbaccess.Selector[[]boardTable.Record],
	log pkgLog.Errorer,
) GETHandler {
	return GETHandler{
		userSelector:           userSelector,
		boardInserter:          boardInserter,
		idValidator:            idValidator,
		boardSelectorRecursive: boardSelectorRecursive,
		teamSelector:           teamSelector,
		userSelectorByTeamID:   userSelectorByTeamID,
		boardSelectorByTeamID:  boardSelectorByTeamID,
		log:                    log,
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

	// Format team ID as string to be used by repository methods.
	teamIDStr := strconv.Itoa(user.TeamID)

	// Select the team from the database that the user is the member/admin of.
	team, err := h.teamSelector.Select(teamIDStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}

	// Get all members of the user's team.
	members, err := h.userSelectorByTeamID.Select(teamIDStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}

	// Get all boards that belong to the user's team.
	boards, err := h.boardSelectorByTeamID.Select(teamIDStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}

	boardID := r.URL.Query().Get("id")
	if boardID == "" {
		if len(boards) > 0 {
			// If the board ID is empty but boards for the user's team were
			// found, set the board ID to the first baord's.
			boardID = strconv.Itoa(boards[0].ID)
		} else if user.IsAdmin {
			// If the board ID was empty and NO boards for the user's team were
			// found, create a board only if the user is an admin.
			if err = h.boardInserter.Insert(
				boardTable.NewInRecord("New Board", user.TeamID),
			); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				h.log.Error(err.Error())
				return
			}
			// Refresh boards and set the board ID to the ID of the first (and
			// only) retrieved.
			boards, err = h.boardSelectorByTeamID.Select(teamIDStr)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				h.log.Error(err.Error())
				return
			}
			boardID = strconv.Itoa(boards[0].ID)
		} else {
			// If the board ID is empty, no boards for the user's team were
			// found, and the user is not admin, return 403 Forbidden.
			w.WriteHeader(http.StatusForbidden)
			return
		}
	} else {
		// Validate board ID if not empty.
		if err = h.idValidator.Validate(boardID); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	if boardID == "" {
		boardID = strconv.Itoa(boards[0].ID)
	}

	// Select recursive board record from the database.
	activeBoard, err := h.boardSelectorRecursive.Select(boardID)
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

	// Build response from data retrieved from the database.
	resp := GETResp{
		User: User{
			Username: user.Username,
			IsAdmin:  user.IsAdmin,
		},
		Team: Team{
			ID:         team.ID,
			InviteCode: team.InviteCode,
		},
		TeamMembers: make([]TeamMember, len(members)),
		Boards:      make([]Board, len(boards)),
		ActiveBoard: ActiveBoard{
			ID:      activeBoard.ID,
			Name:    activeBoard.Name,
			Columns: make([]Column, len(activeBoard.Columns)),
		},
	}

	for i, member := range members {
		resp.TeamMembers[i] = TeamMember{
			Username: member.Username,
			IsAdmin:  member.IsAdmin,
		}
	}

	for i, board := range boards {
		resp.Boards[i] = Board{ID: board.ID, Name: board.Name}
	}

	for i, col := range activeBoard.Columns {
		resp.ActiveBoard.Columns[i].ID = col.ID
		resp.ActiveBoard.Columns[i].Order = col.Order
		resp.ActiveBoard.Columns[i].Tasks = make([]Task, len(col.Tasks))

		for j, task := range col.Tasks {
			resp.ActiveBoard.Columns[i].Tasks[j].ID = task.ID
			resp.ActiveBoard.Columns[i].Tasks[j].Title = task.Title
			var desc string
			if task.Description != nil {
				desc = *task.Description
			}
			resp.ActiveBoard.Columns[i].Tasks[j].Description = desc
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
