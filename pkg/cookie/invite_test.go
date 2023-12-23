//go:build utest

package cookie

import (
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/kxplxn/goteam/pkg/assert"
)

func TestInvite(t *testing.T) {
	key := []byte("signkey")
	teamID := "teamid"

	t.Run("Encode", func(t *testing.T) {
		sut := NewInviteEncoder(key, 1*time.Hour)

		ck, err := sut.Encode(NewInvite(teamID))
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
		sut := NewInviteDecoder(key)

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
				token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZWFtSUQiOiJ0" +
					"ZWFtaWQifQ.1h_fmLJ1ip-Z6kJq9JXYDgGuWDPOcOf8abwCgKtHHcY",
				wantTeamID: "teamid",
				wantErr:    nil,
			},
		} {
			t.Run(c.name, func(t *testing.T) {
				inv, err := sut.Decode(http.Cookie{Value: c.token})

				assert.ErrIs(t.Error, err, c.wantErr)
				assert.Equal(t.Error, inv.TeamID, c.wantTeamID)
			})
		}

	})
}
