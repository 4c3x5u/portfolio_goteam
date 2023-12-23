//go:build utest

package cookie

import (
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/kxplxn/goteam/pkg/assert"
)

func TestState(t *testing.T) {
	key := []byte("signkey")
	boards := []Board{
		NewBoard("boardid", []Column{NewColumn([]Task{
			NewTask("taskid", 2),
		})}),
	}

	t.Run("Encode", func(t *testing.T) {
		sut := NewStateEncoder(key, 1*time.Hour)

		ck, err := sut.Encode(NewState(boards))
		if err != nil {
			t.Fatal(err)
		}

		assert.Nil(t.Fatal, ck.Valid())
		assert.Equal(t.Error, ck.Name, AuthName)
		assert.Equal(t.Error, ck.SameSite, http.SameSiteNoneMode)
		assert.True(t.Error, ck.Secure)
		assert.True(t.Error,
			ck.Expires.UTC().After(time.Now().Add(59*time.Minute).UTC()))
		assert.True(t.Error,
			ck.Expires.UTC().Before(time.Now().Add(61*time.Minute).UTC()))

		claims := jwt.MapClaims{}
		if _, err = jwt.ParseWithClaims(
			ck.Value, &claims, func(token *jwt.Token) (any, error) {
				return key, nil
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

			columns := boardRaw["columns"].([]any)
			col := columns[0].(map[string]any)

			tasks := col["tasks"].([]any)
			ta := tasks[0].(map[string]any)
			taID := ta["id"].(string)
			assert.Equal(t.Error, taID, boards[0].Columns[0].Tasks[0].ID)
			taOrder := int(ta["order"].(float64))
			assert.Equal(t.Error, taOrder, boards[0].Columns[0].Tasks[0].Order)
		}

		assert.True(t.Error,
			int64(claims["exp"].(float64)) >
				time.Now().Add(59*time.Minute).Unix(),
		)
		assert.True(t.Error,
			int64(claims["exp"].(float64)) <
				time.Now().Add(61*time.Minute).Unix(),
		)
	})

	t.Run("Decode", func(t *testing.T) {
		sut := NewStateDecoder(key)

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
				name: "Success",
				token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJib2FyZHMiOlt7" +
					"ImlkIjoiYm9hcmRpZCIsImNvbHVtbnMiOlt7InRhc2tzIjpbeyJpZCI6" +
					"InRhc2tpZCIsIm9yZGVyIjoyfV19XX1dfQ._LZ3QROcAY0n6LbPsqvUF" +
					"ugCD9JQ4CYco00BmrS3Ukc",
				wantBoards: boards,
				wantErr:    nil,
			},
		} {
			t.Run(c.name, func(t *testing.T) {
				st, err := sut.Decode(http.Cookie{Value: c.token})

				assert.ErrIs(t.Fatal, err, c.wantErr)

				for i, b := range c.wantBoards {
					assert.Equal(t.Error, st.Boards[i].ID, b.ID)
					assert.AllEqual(t.Error,
						st.Boards[i].Columns[0].Tasks,
						b.Columns[0].Tasks,
					)
				}
			})
		}

	})
}
