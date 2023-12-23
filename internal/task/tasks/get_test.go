//go:build utest

package tasks

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/cookie"
	"github.com/kxplxn/goteam/pkg/db"
	"github.com/kxplxn/goteam/pkg/db/tasktable"
	pkgLog "github.com/kxplxn/goteam/pkg/log"
)

func TestGetHandler(t *testing.T) {
	authDecoder := &cookie.FakeDecoder[cookie.Auth]{}
	retriever := &db.FakeRetriever[[]tasktable.Task]{}
	log := &pkgLog.FakeErrorer{}
	sut := NewGetHandler(authDecoder, retriever, log)

	someTasks := []tasktable.Task{
		{
			TeamID:       "team1",
			BoardID:      "board1",
			ColumnNumber: 0,
			ID:           "task1",
			Title:        "taskone",
			Description:  "task one description",
			Order:        1,
			Subtasks: []tasktable.Subtask{
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
			Subtasks: []tasktable.Subtask{
				{Title: "subtaskone", IsDone: false},
				{Title: "subtasktwo", IsDone: true},
			},
		},
	}

	for _, c := range []struct {
		name          string
		authToken     string
		errDecodeAuth error
		errRetrieve   error
		tasks         []tasktable.Task
		wantStatus    int
		assertFunc    func(*testing.T, *http.Response, string)
	}{
		{
			name:          "NoAuth",
			authToken:     "",
			errDecodeAuth: nil,
			errRetrieve:   nil,
			tasks:         []tasktable.Task{},
			wantStatus:    http.StatusUnauthorized,
			assertFunc:    func(*testing.T, *http.Response, string) {},
		},
		{
			name:          "InvalidAuth",
			authToken:     "nonempty",
			errDecodeAuth: errors.New("decode auth failed"),
			errRetrieve:   nil,
			tasks:         []tasktable.Task{},
			wantStatus:    http.StatusUnauthorized,
			assertFunc:    func(*testing.T, *http.Response, string) {},
		},
		{
			name:          "ErrRetrieve",
			authToken:     "nonempty",
			errDecodeAuth: nil,
			errRetrieve:   errors.New("retrieve failed"),
			tasks:         []tasktable.Task{},
			wantStatus:    http.StatusInternalServerError,
			assertFunc:    func(*testing.T, *http.Response, string) {},
		},
		{
			name:          "OKSome",
			authToken:     "nonempty",
			errDecodeAuth: nil,
			errRetrieve:   nil,
			tasks:         someTasks,
			wantStatus:    http.StatusOK,
			assertFunc: func(t *testing.T, r *http.Response, _ string) {
				var tasks []tasktable.Task
				err := json.NewDecoder(r.Body).Decode(&tasks)
				if err != nil {
					t.Fatal(err)
				}

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
			authDecoder.Err = c.errDecodeAuth
			retriever.Err = c.errRetrieve
			retriever.Res = c.tasks

			r := httptest.NewRequest(http.MethodGet, "/", nil)
			if c.authToken != "" {
				r.AddCookie(&http.Cookie{
					Name: "auth-token", Value: c.authToken,
				})
			}
			w := httptest.NewRecorder()

			sut.Handle(w, r, "")

			res := w.Result()
			assert.Equal(t.Error, res.StatusCode, c.wantStatus)

			c.assertFunc(t, res, "")
		})
	}
}
