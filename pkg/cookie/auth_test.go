//go:build utest

package cookie

import (
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/kxplxn/goteam/pkg/assert"
)

func TestAuth(t *testing.T) {
	key := []byte("signkey")
	username := "bob123"
	isAdmin := true
	teamID := "teamid"

	t.Run("Encode", func(t *testing.T) {
		dur := 1 * time.Hour
		sut := NewAuthEncoder(key, dur)

		ck, err := sut.Encode(NewAuth(username, isAdmin, teamID))
		assert.Nil(t.Fatal, err)

		assert.Nil(t.Fatal, ck.Valid())
		assert.Equal(t.Error, ck.Name, AuthName)
		assert.Equal(t.Error, ck.SameSite, http.SameSiteNoneMode)
		assert.True(t.Error, ck.Secure)
		assert.True(t.Error,
			ck.Expires.UTC().After(time.Now().Add(59*time.Minute).UTC()))
		assert.True(t.Error,
			ck.Expires.UTC().Before(time.Now().Add(61*time.Minute).UTC()))

		claims := jwt.MapClaims{}
		_, err = jwt.ParseWithClaims(
			ck.Value, &claims, func(token *jwt.Token) (any, error) {
				return key, nil
			},
		)
		assert.Nil(t.Fatal, err)

		assert.Equal(t.Error, claims["username"].(string), username)
		assert.Equal(t.Error, claims["isAdmin"].(bool), isAdmin)
		assert.Equal(t.Error, claims["teamID"].(string), teamID)
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
		sut := NewAuthDecoder(key)

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
				name: "Success",
				token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJib2FyZElEcyI6" +
					"WyJib2FyZDEiLCJib2FyZDIiXSwiaXNBZG1pbiI6dHJ1ZSwidGVhbUlE" +
					"IjoidGVhbWlkIiwidXNlcm5hbWUiOiJib2IxMjMifQ.4uS5-QGd3Gj2I" +
					"5Jm2-dnq-1-_3IqcBepLBdzPjjfRuM",
				wantUsername: username,
				wantIsAdmin:  isAdmin,
				wantTeamID:   teamID,
				wantErr:      nil,
			},
		} {
			t.Run(c.name, func(t *testing.T) {
				auth, err := sut.Decode(http.Cookie{Value: c.token})

				assert.ErrIs(t.Error, err, c.wantErr)
				assert.Equal(t.Error, auth.Username, c.wantUsername)
				assert.Equal(t.Error, auth.IsAdmin, c.wantIsAdmin)
				assert.Equal(t.Error, auth.TeamID, c.wantTeamID)
			})
		}

	})
}
