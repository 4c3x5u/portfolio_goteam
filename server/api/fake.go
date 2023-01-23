package api

import "net/http"

// FakeMethodHandler is a test fake for MethodHandler.
type FakeMethodHandler struct {
	InResponseWriter http.ResponseWriter
	InReq            *http.Request
	InSub            string
}

// Handle implements the MethodHandler interface on FakeMethodHandler.
func (f *FakeMethodHandler) Handle(
	w http.ResponseWriter, r *http.Request, sub string,
) {
	f.InResponseWriter, f.InReq, f.InSub = w, r, sub
}
