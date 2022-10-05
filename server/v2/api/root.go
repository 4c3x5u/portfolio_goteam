package api

import (
	"github.com/kxplxn/goteam/server/v2/log"
	"net/http"
)

func ServeRoot(w http.ResponseWriter, _ *http.Request) {
	relay := log.NewAPILogger(w)
	relay.Msg(
		"app status: OK\n" +
			"available endpoints:\n" +
			"/register",
	)
}
