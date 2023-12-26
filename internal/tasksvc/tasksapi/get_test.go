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
	retriever := &db.FakeRetriever[[]tasktbl.Task]{}
	log := &log.FakeErrorer{}
	sut := NewGetHandler(boardIDValidator, stateDecoder, retriever, log)

	someTasks := []tasktbl.Task{
		{
			TeamID:       "team1",
			BoardID:      "board1",
			ColumnNumber: 0,
			ID:           "task1",
			Title:        "taskone",
			Description:  "task one description",
			Order:        1,
			Subtasks: []tasktbl.Subtask{
				{Title: "subtaskone", IsDone: false},
				{Title: "subtasktwo", IsDone: true},
			},
		},
		{
			TeamID:       "team1",
			BoardID:      "board2",
			ColumnNumber: 2,
			ID:           "task1",
			Title:        "taskone",
			Description:  "task one description",
			Order:        3,
			Subtasks: []tasktbl.Subtask{
				{Title: "subtaskone", IsDone: false},
				{Title: "subtasktwo", IsDone: true},
			},
		},
	}

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
			errDecodeState:     errors.New("decode auth failed"),
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
				ID: "adksjh",
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
				ID: "adksjh",
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
				ID: "adksjh",
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
					assert.Equal(t.Error, gotTask.TeamID, someTasks[i].TeamID)
					assert.Equal(t.Error, gotTask.BoardID, someTasks[i].BoardID)
					assert.Equal(t.Error,
						gotTask.ColumnNumber, someTasks[i].ColumnNumber,
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
							gotSubtask.Title, someTasks[i].Subtasks[j].Title,
						)
						assert.Equal(t.Error,
							gotSubtask.IsDone, someTasks[i].Subtasks[j].IsDone,
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
			retriever.Err = c.errRetrieve
			retriever.Res = c.tasks
			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				http.MethodGet, "/?boardID=adksjh", nil,
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
}
