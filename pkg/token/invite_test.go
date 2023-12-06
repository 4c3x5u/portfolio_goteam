//go:build utest

package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/kxplxn/goteam/pkg/assert"
)

func TestInvite(t *testing.T) {
	signKey = []byte("signkey")
	teamID := "teamid"
	sut := NewInvite()

	t.Run("Encode", func(t *testing.T) {
		sut.TeamID = teamID
		expiry := time.Now().Add(1 * time.Hour)

		token, err := sut.Encode(expiry)
		if err != nil {
			t.Fatal(err)
		}

		claims := jwt.MapClaims{}
		if _, err = jwt.ParseWithClaims(
			token,
			&claims,
			func(token *jwt.Token) (any, error) { return signKey, nil },
		); err != nil {
			t.Error(err)
		}

		assert.Equal(t.Error, claims["teamID"].(string), teamID)
		assert.Equal(t.Error, claims["exp"].(float64), float64(expiry.Unix()))
	})

	t.Run("Decode", func(t *testing.T) {
		for _, c := range []struct {
			name       string
			token      string
			wantTeamID string
			wantErr    error
		}{
			{
				name: "InvalidSignature",
				token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDE5" +
					"MDE0NzMsInRlYW1JRCI6ImFza2RqZmFza2RmamFoIn0.g0jCuok1he1o" +
					"puHRGfVmvuGtpwfWlIBbnRK64qgLsx4",
				wantTeamID: "",
				wantErr:    jwt.ErrSignatureInvalid,
			},
			{
				name: "TokenMalformed",
				token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCb2IyMSJ9.k6QDVjyaHx" +
					"PYixeoQBLixC5c79VK-WZ_kD9u4fjX_Ks",
				wantTeamID: "",
				wantErr:    jwt.ErrTokenMalformed,
			},
			{
				name: "Expired",
				token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDE4" +
					"OTgwOTYsInRlYW1JRCI6InRlYW1pZCJ9.MAORMCFqzNrLnY4l_wrPA86" +
					"K9w6W9pzH_4b6iNHr1SE",
				wantTeamID: "",
				wantErr:    jwt.ErrTokenExpired,
			},
			{
				name: "Success",
				token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDE5" +
					"MDE4MTksInRlYW1JRCI6InRlYW1pZCJ9.J886zVP8YWquNHf03A9QK3m" +
					"rC389Fd9NiDzLVYtrHC4",
				wantTeamID: "teamid",
				wantErr:    nil,
			},
		} {
			t.Run(c.name, func(t *testing.T) {
				err := sut.Decode(c.token)

				assert.ErrIs(t.Error, err, c.wantErr)
				assert.Equal(t.Error, sut.TeamID, c.wantTeamID)
			})
		}

	})
}
