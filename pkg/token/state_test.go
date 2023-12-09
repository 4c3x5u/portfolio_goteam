//go:build utest

package token

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/kxplxn/goteam/pkg/assert"
)

func TestState(t *testing.T) {
	if err := os.Setenv(keyName, "signkey"); err != nil {
		t.Fatal("failed to set key env var", err)
	}
	boards := []Board{
		NewBoard("boardid", []Column{NewColumn(1)}),
	}

	t.Run("Encode", func(t *testing.T) {
		// arrange
		expiry := time.Now().Add(1 * time.Hour)

		// act
		token, err := EncodeState(expiry, NewState(boards))
		if err != nil {
			t.Fatal(err)
		}

		// assert
		claims := jwt.MapClaims{}
		if _, err = jwt.ParseWithClaims(token, &claims,
			func(token *jwt.Token) (any, error) {
				return []byte(os.Getenv(keyName)), nil
			},
		); err != nil {
			t.Error(err)
		}

		boardsRaw := claims["boards"].([]any)

		for i := range boards {
			b := boardsRaw[i]

			boardRaw := b.(map[string]any)

			id := boardRaw["id"].(string)
			assert.Equal(t.Error, id, boards[0].ID)

			columnsRaw := boardRaw["columns"].([]any)

			col := columnsRaw[0].(map[string]any)
			taskCount := col["taskCount"].(float64)
			assert.Equal(t.Error,
				taskCount, float64(boards[0].Columns[0].TaskCount),
			)
		}
	})

	t.Run("Decode", func(t *testing.T) {
		for _, c := range []struct {
			name       string
			token      string
			wantBoards []Board
			wantErr    error
		}{
			{
				name: "InvalidSignature",
				token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDE5" +
					"MDE0NzMsInRlYW1JRCI6ImFza2RqZmFza2RmamFoIn0.g0jCuok1he1o" +
					"puHRGfVmvuGtpwfWlIBbnRK64qgLsx4",
				wantBoards: []Board{},
				wantErr:    jwt.ErrSignatureInvalid,
			},
			{
				name: "TokenMalformed",
				token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCb2IyMSJ9.k6QDVjyaHx" +
					"PYixeoQBLixC5c79VK-WZ_kD9u4fjX_Ks",
				wantBoards: []Board{},
				wantErr:    jwt.ErrTokenMalformed,
			},
			{
				name: "Expired",
				token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDE4" +
					"OTgwOTYsInRlYW1JRCI6InRlYW1pZCJ9.MAORMCFqzNrLnY4l_wrPA86" +
					"K9w6W9pzH_4b6iNHr1SE",
				wantBoards: []Board{},
				wantErr:    jwt.ErrTokenExpired,
			},
			{
				name:       "Success",
				token:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJib2FyZHMiOlt7ImlkIjoiYm9hcmRpZCIsImNvbHVtbnMiOlt7InRhc2tDb3VudCI6MX1dfV19.vVhDRAKKHD9FAu9FIV__N74HYm3UBKhD_4Z_bJujpDw",
				wantBoards: boards,
				wantErr:    nil,
			},
		} {
			t.Run(c.name, func(t *testing.T) {
				st, err := DecodeState(c.token)

				assert.ErrIs(t.Fatal, err, c.wantErr)

				for i, b := range c.wantBoards {
					assert.Equal(t.Error, st.Boards[i].ID, b.ID)
					assert.AllEqual(t.Error, st.Boards[i].Columns, b.Columns)
				}
			})
		}

	})
}
