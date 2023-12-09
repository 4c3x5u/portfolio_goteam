//go:build utest

package token

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/kxplxn/goteam/pkg/assert"
)

func TestAuth(t *testing.T) {
	if err := os.Setenv(keyName, "signkey"); err != nil {
		t.Fatal("failed to set key env var", err)
	}
	username := "bob123"
	isAdmin := true
	teamID := "teamid"

	t.Run("Encode", func(t *testing.T) {
		expiry := time.Now().Add(1 * time.Hour)

		token, err := EncodeAuth(expiry, NewAuth(username, isAdmin, teamID))
		if err != nil {
			t.Fatal(err)
		}

		claims := jwt.MapClaims{}
		if _, err = jwt.ParseWithClaims(
			token,
			&claims,
			func(token *jwt.Token) (any, error) {
				return []byte(os.Getenv(keyName)), nil
			},
		); err != nil {
			t.Error(err)
		}

		assert.Equal(t.Error, claims["username"].(string), username)
		assert.Equal(t.Error, claims["isAdmin"].(bool), isAdmin)
		assert.Equal(t.Error, claims["teamID"].(string), teamID)
	})

	t.Run("Decode", func(t *testing.T) {
		for _, c := range []struct {
			name         string
			token        string
			wantUsername string
			wantIsAdmin  bool
			wantTeamID   string
			wantErr      error
		}{
			{
				name: "InvalidSignature",
				token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDE5" +
					"MDE0NzMsInRlYW1JRCI6ImFza2RqZmFza2RmamFoIn0.g0jCuok1he1o" +
					"puHRGfVmvuGtpwfWlIBbnRK64qgLsx4",
				wantUsername: "",
				wantIsAdmin:  false,
				wantTeamID:   "",
				wantErr:      jwt.ErrSignatureInvalid,
			},
			{
				name: "TokenMalformed",
				token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCb2IyMSJ9.k6QDVjyaHx" +
					"PYixeoQBLixC5c79VK-WZ_kD9u4fjX_Ks",
				wantUsername: "",
				wantIsAdmin:  false,
				wantTeamID:   "",
				wantErr:      jwt.ErrTokenMalformed,
			},
			{
				name: "Expired",
				token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDE4" +
					"OTgwOTYsInRlYW1JRCI6InRlYW1pZCJ9.MAORMCFqzNrLnY4l_wrPA86" +
					"K9w6W9pzH_4b6iNHr1SE",
				wantUsername: "",
				wantIsAdmin:  false,
				wantTeamID:   "",
				wantErr:      jwt.ErrTokenExpired,
			},
			{
				name:         "Success",
				token:        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJib2FyZElEcyI6WyJib2FyZDEiLCJib2FyZDIiXSwiaXNBZG1pbiI6dHJ1ZSwidGVhbUlEIjoidGVhbWlkIiwidXNlcm5hbWUiOiJib2IxMjMifQ.4uS5-QGd3Gj2I5Jm2-dnq-1-_3IqcBepLBdzPjjfRuM",
				wantUsername: username,
				wantIsAdmin:  isAdmin,
				wantTeamID:   teamID,
				wantErr:      nil,
			},
		} {
			t.Run(c.name, func(t *testing.T) {
				inv, err := DecodeAuth(c.token)

				assert.ErrIs(t.Error, err, c.wantErr)
				assert.Equal(t.Error, inv.Username, c.wantUsername)
				assert.Equal(t.Error, inv.IsAdmin, c.wantIsAdmin)
				assert.Equal(t.Error, inv.TeamID, c.wantTeamID)
			})
		}

	})
}
