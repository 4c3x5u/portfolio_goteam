package api

import (
	"net/http"

	"github.com/kxplxn/goteam/server/v2/relay"
)

type HandlerRoot struct {
}

func (h *HandlerRoot) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	relay.New(w).APIMsg(
		"app status: OK\n" +
			"   available endpoints: " +
			"/register",
	)
}
