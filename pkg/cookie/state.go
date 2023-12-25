package cookie

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// StateName is the name of the state token.
const StateName = "state-token"

// State defines the body of the state token.
type State struct{ Boards []Board }

// NewState creates and returns a new State.
func NewState(boards []Board) State {
	return State{Boards: boards}
}

// Board defines Board data that the state token contains.
type Board struct {
	ID      string   `json:"id"`
	Columns []Column `json:"columns"`
}

// NewBoard creates and returns a new Board.
func NewBoard(id string, columns []Column) Board {
	return Board{ID: id, Columns: columns}
}

// Column defines an element in the column data that a Board contains.
type Column struct {
	Tasks []Task `json:"tasks"`
}

// NewColumn creates and returns a new Column
func NewColumn(tasks []Task) Column { return Column{Tasks: tasks} }

// Task defines an element in the tasks that a Column contains.
type Task struct {
	ID    string `json:"id"`
	Order int    `json:"order"`
}

// NewTask creates and returns a new Task.
func NewTask(id string, order int) Task { return Task{ID: id, Order: order} }

// StateEncoder defines a type that can be used to encode an state token.
type StateEncoder struct {
	key []byte
	dur time.Duration
}

// NewStateEncoder creates and returns a new StateEncoder.
func NewStateEncoder(jwtKey []byte, duration time.Duration) StateEncoder {
	return StateEncoder{key: jwtKey, dur: duration}
}

// Encode encodes an State into a JWT string.
func (e StateEncoder) Encode(state State) (http.Cookie, error) {
	exp := time.Now().Add(e.dur)

	tk, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"boards": state.Boards,
		"exp":    exp.Unix(),
	}).SignedString(e.key)
	if err != nil {
		return http.Cookie{}, err
	}

	return http.Cookie{
		Name:     AuthName,
		Value:    tk,
		Expires:  exp.UTC(),
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}, nil
}

// StateDecoder defines a type that can be used to decode an state token.
type StateDecoder struct{ key []byte }

// NewStateDecoder creates and returns a new StateDecoder.
func NewStateDecoder(jwtKey []byte) StateDecoder {
	return StateDecoder{key: jwtKey}
}

// Decode validates and decodes a raw JWT string into an State.
func (d StateDecoder) Decode(ck http.Cookie) (State, error) {
	if ck.Value == "" {
		return State{}, ErrInvalid
	}

	claims := jwt.MapClaims{}
	if _, err := jwt.ParseWithClaims(
		ck.Value, &claims, func(token *jwt.Token) (any, error) {
			return d.key, nil
		},
	); err != nil {
		return State{}, err
	}

	boardsRaw, ok := claims["boards"].([]any)
	if !ok {
		return State{}, ErrInvalid
	}

	var boards []Board
	for _, b := range boardsRaw {
		boardRaw, ok := b.(map[string]any)
		if !ok {
			return State{}, ErrInvalid
		}

		id, ok := boardRaw["id"].(string)
		if !ok {
			return State{}, ErrInvalid
		}

		columnsRaw, ok := boardRaw["columns"].([]any)
		if !ok {
			return State{}, ErrInvalid
		}

		var columns []Column
		for _, c := range columnsRaw {
			colRaw, ok := c.(map[string]any)
			if !ok {
				return State{}, ErrInvalid
			}

			tasksRaw, ok := colRaw["tasks"].([]any)
			if !ok {
				return State{}, ErrInvalid
			}

			var tasks []Task
			for _, t := range tasksRaw {
				tRaw, ok := t.(map[string]any)
				if !ok {
					return State{}, ErrInvalid
				}

				tID, ok := tRaw["id"].(string)
				if !ok {
					return State{}, ErrInvalid
				}

				tOrder, ok := tRaw["order"].(float64)
				if !ok {
					return State{}, ErrInvalid
				}

				tasks = append(tasks, NewTask(tID, int(tOrder)))
			}

			columns = append(columns, NewColumn(tasks))
		}

		boards = append(boards, NewBoard(id, columns))
	}

	return NewState(boards), nil
}
