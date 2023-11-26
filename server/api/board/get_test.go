//go:build utest

package board

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kxplxn/goteam/server/api"
	"github.com/kxplxn/goteam/server/assert"
	boardTable "github.com/kxplxn/goteam/server/dbaccess/board"
	teamTable "github.com/kxplxn/goteam/server/dbaccess/team"
	userTable "github.com/kxplxn/goteam/server/dbaccess/user"
	pkgLog "github.com/kxplxn/goteam/server/log"
)

// TestGETHandler tests the Handle method of GETHandler to assert that it
// behaves correctly in all possible scenarios.
func TestGETHandler(t *testing.T) {
	userSelector := &userTable.FakeSelector{}
	idValidator := &api.FakeStringValidator{}
	boardSelector := &boardTable.FakeRecursiveSelector{}
	teamSelector := &teamTable.FakeSelector{}
	log := &pkgLog.FakeErrorer{}

	sut := NewGETHandler(
		userSelector, idValidator, boardSelector, teamSelector, log,
	)

	for _, c := range []struct {
		name             string
		user             userTable.Record
		userSelectorErr  error
		idValidatorErr   error
		board            boardTable.RecursiveRecord
		boardSelectorErr error
		team             teamTable.Record
		teamSelectorErr  error
		wantStatusCode   int
		assertFunc       func(*testing.T, *http.Response, string)
	}{
		{
			name:             "UserIsNotRecognised",
			user:             userTable.Record{},
			userSelectorErr:  sql.ErrNoRows,
			idValidatorErr:   nil,
			board:            boardTable.RecursiveRecord{},
			boardSelectorErr: nil,
			team:             teamTable.Record{},
			teamSelectorErr:  nil,
			wantStatusCode:   http.StatusUnauthorized,
			assertFunc:       assert.OnLoggedErr(sql.ErrNoRows.Error()),
		},
		{
			name:             "UserSelectorErr",
			user:             userTable.Record{},
			userSelectorErr:  sql.ErrConnDone,
			idValidatorErr:   nil,
			board:            boardTable.RecursiveRecord{},
			boardSelectorErr: nil,
			team:             teamTable.Record{},
			teamSelectorErr:  nil,
			wantStatusCode:   http.StatusInternalServerError,
			assertFunc:       assert.OnLoggedErr(sql.ErrConnDone.Error()),
		},
		{
			name:             "InvalidID",
			user:             userTable.Record{},
			userSelectorErr:  nil,
			idValidatorErr:   errors.New("error invalid id"),
			board:            boardTable.RecursiveRecord{},
			boardSelectorErr: nil,
			team:             teamTable.Record{},
			teamSelectorErr:  nil,
			wantStatusCode:   http.StatusBadRequest,
			assertFunc:       func(_ *testing.T, _ *http.Response, _ string) {},
		},
		{
			name:             "BoardNotFound",
			user:             userTable.Record{},
			userSelectorErr:  nil,
			idValidatorErr:   nil,
			board:            boardTable.RecursiveRecord{},
			boardSelectorErr: sql.ErrNoRows,
			team:             teamTable.Record{},
			teamSelectorErr:  nil,
			wantStatusCode:   http.StatusNotFound,
			assertFunc:       func(_ *testing.T, _ *http.Response, _ string) {},
		},
		{
			name:             "BoardSelectorErr",
			user:             userTable.Record{},
			userSelectorErr:  nil,
			idValidatorErr:   nil,
			board:            boardTable.RecursiveRecord{},
			boardSelectorErr: sql.ErrConnDone,
			team:             teamTable.Record{},
			teamSelectorErr:  nil,
			wantStatusCode:   http.StatusInternalServerError,
			assertFunc:       assert.OnLoggedErr(sql.ErrConnDone.Error()),
		},
		{
			name:             "BoardWrongTeam",
			user:             userTable.Record{TeamID: 1},
			userSelectorErr:  nil,
			idValidatorErr:   nil,
			board:            boardTable.RecursiveRecord{TeamID: 2},
			boardSelectorErr: nil,
			team:             teamTable.Record{},
			teamSelectorErr:  nil,
			wantStatusCode:   http.StatusForbidden,
			assertFunc:       func(_ *testing.T, _ *http.Response, _ string) {},
		},
		{
			name:             "TeamSelectorErr",
			user:             userTable.Record{},
			userSelectorErr:  nil,
			idValidatorErr:   nil,
			board:            boardTable.RecursiveRecord{},
			boardSelectorErr: nil,
			team:             teamTable.Record{},
			teamSelectorErr:  sql.ErrNoRows,
			wantStatusCode:   http.StatusInternalServerError,
			assertFunc:       assert.OnLoggedErr(sql.ErrNoRows.Error()),
		},
		{
			name:            "OK",
			user:            userTable.Record{},
			userSelectorErr: nil,
			idValidatorErr:  nil,
			board: func() boardTable.RecursiveRecord {
				task1Desc := "task1Desc"
				task2Desc := "task2Desc"
				return boardTable.RecursiveRecord{
					ID: 2, Name: "Active", Columns: []boardTable.Column{
						{ID: 3, Order: 1, Tasks: []boardTable.Task{}},
						{ID: 4, Order: 2, Tasks: []boardTable.Task{
							{
								ID:          5,
								Title:       "task1title",
								Description: &task1Desc,
								Order:       3,
								Subtasks: []boardTable.Subtask{
									{
										ID:     5,
										Title:  "subtask1",
										Order:  4,
										IsDone: true,
									},
									{
										ID:     6,
										Title:  "subtask2",
										Order:  5,
										IsDone: false,
									},
								},
							},
							{
								ID:          7,
								Title:       "task2title",
								Description: &task2Desc,
								Order:       6,
								Subtasks:    []boardTable.Subtask{},
							},
						}},
					},
				}
			}(),
			boardSelectorErr: nil,
			team:             teamTable.Record{ID: 1, InviteCode: "InvCode"},
			teamSelectorErr:  nil,
			wantStatusCode:   http.StatusOK,
			assertFunc: func(t *testing.T, r *http.Response, _ string) {
				var resp GETResp
				if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
					t.Fatal(err)
				}

				if err := assert.Equal("bob123", resp.Username); err != nil {
					t.Error(err)
				}
				if err := assert.Equal(1, resp.Team.ID); err != nil {
					t.Error(err)
				}
				if err := assert.Equal(
					"InvCode", resp.Team.InviteCode,
				); err != nil {
					t.Error(err)
				}

				// TODO: assert on all team members
				if err := assert.Equal(0, len(resp.TeamMembers)); err != nil {
					t.Error(err)
				}

				// TODO: assert on all boards
				if err := assert.Equal(0, len(resp.Boards)); err != nil {
					t.Error(err)
				}

				if err := assert.Equal(2, resp.ActiveBoard.ID); err != nil {
					t.Error(err)
				}
				if err := assert.Equal(
					"Active", resp.ActiveBoard.Name,
				); err != nil {
					t.Error(err)
				}
				if err := assert.Equal(
					2, len(resp.ActiveBoard.Columns),
				); err != nil {
					t.Error(err)
				}
				for i, col := range resp.ActiveBoard.Columns {
					wantCol := boardSelector.Rec.Columns[i]
					if err := assert.Equal(wantCol.ID, col.ID); err != nil {
						t.Error(err)
					}
					if err := assert.Equal(
						wantCol.Order, col.Order,
					); err != nil {
						t.Error(err)
					}
					if err := assert.Equal(
						len(wantCol.Tasks), len(col.Tasks),
					); err != nil {
						t.Error(err)
					}
					for j, task := range col.Tasks {
						wantTask := wantCol.Tasks[j]
						if err := assert.Equal(
							wantTask.ID,
							task.ID,
						); err != nil {
							t.Error(err)
						}
						if err := assert.Equal(
							wantTask.Title,
							task.Title,
						); err != nil {
							t.Error(err)
						}
						if err := assert.Equal(
							*wantTask.Description,
							task.Description,
						); err != nil {
							t.Error(err)
						}

						if err := assert.Equal(
							len(wantTask.Subtasks),
							len(task.Subtasks),
						); err != nil {
							t.Error(err)
						}
						for k, subtask := range task.Subtasks {
							wantSubtask := wantTask.Subtasks[k]
							if err := assert.Equal(
								wantSubtask.ID,
								subtask.ID,
							); err != nil {
								t.Error(err)
							}
							if err := assert.Equal(
								wantSubtask.Title,
								subtask.Title,
							); err != nil {
								t.Error(err)
							}
							if err := assert.Equal(
								wantSubtask.Order,
								subtask.Order,
							); err != nil {
								t.Error(err)
							}
							if err := assert.Equal(
								wantSubtask.IsDone,
								subtask.IsDone,
							); err != nil {
								t.Error(err)
							}
						}
					}
				}
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			userSelector.User = c.user
			userSelector.Err = c.userSelectorErr
			idValidator.Err = c.idValidatorErr
			boardSelector.Rec = c.board
			boardSelector.Err = c.boardSelectorErr
			teamSelector.Rec = c.team
			teamSelector.Err = c.teamSelectorErr

			r, err := http.NewRequest(http.MethodGet, "?boardID=1", nil)
			if err != nil {
				t.Fatal(err)
			}
			w := httptest.NewRecorder()

			sut.Handle(w, r, "bob123")

			res := w.Result()
			if err := assert.Equal(
				c.wantStatusCode, res.StatusCode,
			); err != nil {
				t.Error(err)
			}
			c.assertFunc(t, res, log.InMessage)
		})
	}
}
