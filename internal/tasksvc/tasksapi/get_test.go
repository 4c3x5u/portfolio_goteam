//go:build utest

package tasksapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/cookie"
	"github.com/kxplxn/goteam/pkg/db"
	"github.com/kxplxn/goteam/pkg/db/tasktbl"
	"github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/pkg/validator"
)

func TestGetHandler(t *testing.T) {
	boardIDValidator := &validator.FakeString{}
	stateDecoder := &cookie.FakeDecoder[cookie.State]{}
	retrieverByBoard := &db.FakeRetriever[[]tasktbl.Task]{}
	authDecoder := &cookie.FakeDecoder[cookie.Auth]{}
	retrieverByTeam := &db.FakeRetriever[[]tasktbl.Task]{}
	log := &log.FakeErrorer{}
	sut := NewGetHandler(
		boardIDValidator,
		stateDecoder,
		retrieverByBoard,
		authDecoder,
		retrieverByTeam,
		log,
	)

	someTasks := []tasktbl.Task{
		{
			TeamID:      "team1",
			BoardID:     "board1",
			ColNo:       0,
			ID:          "task1",
			Title:       "taskone",
			Description: "task one description",
			Order:       1,
			Subtasks: []tasktbl.Subtask{
				{Title: "subtaskone", IsDone: false},
				{Title: "subtasktwo", IsDone: false},
			},
		},
		{
			TeamID:      "team1",
			BoardID:     "board1",
			ColNo:       2,
			ID:          "task2",
			Title:       "tasktwo",
			Description: "task two description",
			Order:       2,
			Subtasks: []tasktbl.Subtask{
				{Title: "subtaskthree", IsDone: true},
				{Title: "subtaskfour", IsDone: false},
			},
		},
		{
			TeamID:      "team1",
			BoardID:     "board2",
			ColNo:       0,
			ID:          "task3",
			Title:       "taskthree",
			Description: "task three description",
			Order:       3,
			Subtasks: []tasktbl.Subtask{
				{Title: "subtaskfive", IsDone: true},
				{Title: "subtasksix", IsDone: true},
			},
		},
	}

	t.Run("WithBoardID", func(t *testing.T) {
		for _, c := range []struct {
			name               string
			errValidateBoardID error
			stateToken         string
			errDecodeState     error
			stateDecoded       cookie.State
			errRetrieve        error
			tasks              []tasktbl.Task
			wantStatus         int
			assertFunc         func(*testing.T, *http.Response, []any)
		}{
			{
				name:               "InvalidBoardID",
				errValidateBoardID: errors.New("validate board ID failed"),
				stateToken:         "",
				errDecodeState:     nil,
				stateDecoded:       cookie.State{},
				errRetrieve:        nil,
				tasks:              []tasktbl.Task{},
				wantStatus:         http.StatusBadRequest,
				assertFunc:         func(*testing.T, *http.Response, []any) {},
			},
			{
				name:               "NoState",
				errValidateBoardID: nil,
				stateToken:         "",
				errDecodeState:     nil,
				stateDecoded:       cookie.State{},
				errRetrieve:        nil,
				tasks:              []tasktbl.Task{},
				wantStatus:         http.StatusUnauthorized,
				assertFunc:         func(*testing.T, *http.Response, []any) {},
			},
			{
				name:               "InvalidState",
				errValidateBoardID: nil,
				stateToken:         "nonempty",
				errDecodeState:     errors.New("decode state failed"),
				stateDecoded:       cookie.State{},
				errRetrieve:        nil,
				tasks:              []tasktbl.Task{},
				wantStatus:         http.StatusUnauthorized,
				assertFunc:         func(*testing.T, *http.Response, []any) {},
			},
			{
				name:               "NoAccess",
				errValidateBoardID: nil,
				stateToken:         "nonempty",
				errDecodeState:     nil,
				stateDecoded:       cookie.State{},
				errRetrieve:        nil,
				tasks:              []tasktbl.Task{},
				wantStatus:         http.StatusUnauthorized,
				assertFunc:         func(*testing.T, *http.Response, []any) {},
			},
			{
				name:               "ErrRetrieve",
				errValidateBoardID: nil,
				stateToken:         "nonempty",
				errDecodeState:     nil,
				stateDecoded: cookie.State{Boards: []cookie.Board{{
					ID: "nonempty",
				}}},
				errRetrieve: errors.New("retrieve failed"),
				tasks:       []tasktbl.Task{},
				wantStatus:  http.StatusInternalServerError,
				assertFunc:  func(*testing.T, *http.Response, []any) {},
			},
			{
				name:               "OKNone",
				errValidateBoardID: nil,
				stateToken:         "nonempty",
				errDecodeState:     nil,
				stateDecoded: cookie.State{Boards: []cookie.Board{{
					ID: "nonempty",
				}}},
				errRetrieve: nil,
				tasks:       []tasktbl.Task{},
				wantStatus:  http.StatusOK,
				assertFunc: func(t *testing.T, resp *http.Response, _ []any) {
					var tasks []tasktbl.Task
					err := json.NewDecoder(resp.Body).Decode(&tasks)
					assert.Nil(t.Fatal, err)

					assert.Equal(t.Error, len(tasks), 0)
				},
			},
			{
				name:               "OKSome",
				errValidateBoardID: nil,
				stateToken:         "nonempty",
				errDecodeState:     nil,
				stateDecoded: cookie.State{Boards: []cookie.Board{{
					ID: "nonempty",
				}}},
				errRetrieve: nil,
				tasks:       someTasks,
				wantStatus:  http.StatusOK,
				assertFunc: func(t *testing.T, resp *http.Response, _ []any) {
					var tasks []tasktbl.Task
					err := json.NewDecoder(resp.Body).Decode(&tasks)
					assert.Nil(t.Fatal, err)

					assert.Equal(t.Error, len(tasks), len(someTasks))
					for i, gotTask := range tasks {
						assert.Equal(t.Error,
							gotTask.TeamID, someTasks[i].TeamID,
						)
						assert.Equal(t.Error,
							gotTask.BoardID, someTasks[i].BoardID,
						)
						assert.Equal(t.Error,
							gotTask.ColNo, someTasks[i].ColNo,
						)
						assert.Equal(t.Error, gotTask.ID, someTasks[i].ID)
						assert.Equal(t.Error, gotTask.Title, someTasks[i].Title)
						assert.Equal(t.Error,
							gotTask.Description, someTasks[i].Description,
						)
						assert.Equal(t.Error, gotTask.Order, someTasks[i].Order)

						assert.Equal(t.Error,
							len(gotTask.Subtasks), len(someTasks[i].Subtasks),
						)
						for j, gotSubtask := range gotTask.Subtasks {
							assert.Equal(t.Error,
								gotSubtask.Title,
								someTasks[i].Subtasks[j].Title,
							)
							assert.Equal(t.Error,
								gotSubtask.IsDone,
								someTasks[i].Subtasks[j].IsDone,
							)
						}
					}
				},
			},
		} {
			t.Run(c.name, func(t *testing.T) {
				boardIDValidator.Err = c.errValidateBoardID
				stateDecoder.Res = c.stateDecoded
				stateDecoder.Err = c.errDecodeState
				retrieverByBoard.Err = c.errRetrieve
				retrieverByBoard.Res = c.tasks
				w := httptest.NewRecorder()
				r := httptest.NewRequest(
					http.MethodGet, "/?boardID=nonempty", nil,
				)
				if c.stateToken != "" {
					r.AddCookie(&http.Cookie{
						Name: "state-token", Value: c.stateToken,
					})
				}

				sut.Handle(w, r, "")

				resp := w.Result()
				assert.Equal(t.Error, resp.StatusCode, c.wantStatus)
				c.assertFunc(t, resp, log.Args)
			})
		}
	})

	t.Run("WithoutBoardID", func(t *testing.T) {
		for _, c := range []struct {
			name               string
			errValidateBoardID error
			authToken          string
			errDecodeAuth      error
			authDecoded        cookie.Auth
			errRetrieve        error
			tasks              []tasktbl.Task
			wantStatus         int
			assertFunc         func(*testing.T, *http.Response, []any)
		}{
			{
				name:               "NoAuth",
				errValidateBoardID: nil,
				authToken:          "",
				errDecodeAuth:      nil,
				authDecoded:        cookie.Auth{},
				errRetrieve:        nil,
				tasks:              []tasktbl.Task{},
				wantStatus:         http.StatusUnauthorized,
				assertFunc:         func(*testing.T, *http.Response, []any) {},
			},
			{
				name:               "InvalidAuth",
				errValidateBoardID: nil,
				authToken:          "nonempty",
				errDecodeAuth:      errors.New("decode auth failed"),
				authDecoded:        cookie.Auth{},
				errRetrieve:        nil,
				tasks:              []tasktbl.Task{},
				wantStatus:         http.StatusUnauthorized,
				assertFunc:         func(*testing.T, *http.Response, []any) {},
			},
			{
				name:               "ErrRetrieve",
				errValidateBoardID: nil,
				authToken:          "nonempty",
				errDecodeAuth:      nil,
				authDecoded:        cookie.Auth{TeamID: "team1"},
				errRetrieve:        errors.New("retrieve failed"),
				tasks:              []tasktbl.Task{},
				wantStatus:         http.StatusInternalServerError,
				assertFunc:         func(*testing.T, *http.Response, []any) {},
			},
			{
				name:               "OKNone",
				errValidateBoardID: nil,
				authToken:          "nonempty",
				errDecodeAuth:      nil,
				authDecoded:        cookie.Auth{TeamID: "team1"},
				errRetrieve:        nil,
				tasks:              []tasktbl.Task{},
				wantStatus:         http.StatusOK,
				assertFunc: func(t *testing.T, resp *http.Response, _ []any) {
					var tasks []tasktbl.Task
					err := json.NewDecoder(resp.Body).Decode(&tasks)
					assert.Nil(t.Fatal, err)

					assert.Equal(t.Error, len(tasks), 0)
				},
			},
			{
				name:               "OKSome",
				errValidateBoardID: nil,
				authToken:          "nonempty",
				errDecodeAuth:      nil,
				authDecoded:        cookie.Auth{TeamID: "team1"},
				errRetrieve:        nil,
				tasks:              someTasks,
				wantStatus:         http.StatusOK,
				assertFunc: func(t *testing.T, resp *http.Response, _ []any) {
					var tasks []tasktbl.Task
					err := json.NewDecoder(resp.Body).Decode(&tasks)
					assert.Nil(t.Fatal, err)

					// only the first two tasks share the same board ID,
					// therefore only the first two tasks should be returned
					wantTasks := someTasks[:2]

					assert.Equal(t.Error, len(tasks), len(wantTasks))
					for i, gotTask := range tasks {
						assert.Equal(t.Error,
							gotTask.TeamID, wantTasks[i].TeamID,
						)
						assert.Equal(t.Error,
							gotTask.BoardID, wantTasks[i].BoardID,
						)
						assert.Equal(t.Error,
							gotTask.ColNo, wantTasks[i].ColNo,
						)
						assert.Equal(t.Error, gotTask.ID, wantTasks[i].ID)
						assert.Equal(t.Error, gotTask.Title, wantTasks[i].Title)
						assert.Equal(t.Error,
							gotTask.Description, wantTasks[i].Description,
						)
						assert.Equal(t.Error, gotTask.Order, wantTasks[i].Order)

						assert.Equal(t.Error,
							len(gotTask.Subtasks), len(wantTasks[i].Subtasks),
						)
						for j, gotSubtask := range gotTask.Subtasks {
							assert.Equal(t.Error,
								gotSubtask.Title,
								wantTasks[i].Subtasks[j].Title,
							)
							assert.Equal(t.Error,
								gotSubtask.IsDone,
								wantTasks[i].Subtasks[j].IsDone,
							)
						}
					}
				},
			},
		} {
			t.Run(c.name, func(t *testing.T) {
				authDecoder.Res = c.authDecoded
				authDecoder.Err = c.errDecodeAuth
				retrieverByTeam.Err = c.errRetrieve
				retrieverByTeam.Res = c.tasks
				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodGet, "/", nil)
				if c.authToken != "" {
					r.AddCookie(&http.Cookie{
						Name: "auth-token", Value: c.authToken,
					})
				}

				sut.Handle(w, r, "")

				resp := w.Result()
				assert.Equal(t.Error, resp.StatusCode, c.wantStatus)
				c.assertFunc(t, resp, log.Args)
			})
		}
	})
}
