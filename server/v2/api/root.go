package api

import (
	"net/http"

	"github.com/kxplxn/goteam/server/v2/relay"
)

func HandleRoot(w http.ResponseWriter, _ *http.Request) {
	relay.APIMsg(w, "app status: OK\navailable endpoints: `/register`")
}
