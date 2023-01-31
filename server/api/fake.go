package api

import "net/http"

// FakeMethodHandler is a test fake for MethodHandler.
type FakeMethodHandler struct {
	InResponseWriter http.ResponseWriter
	InReq            *http.Request
	InSub            string
}

// Handle implements the MethodHandler interface on FakeMethodHandler. It
// assigns the parameters passed into it to their corresponding In... fields on
// the fake instance and returns its Out.. fields as per function signature.
func (f *FakeMethodHandler) Handle(
	w http.ResponseWriter, r *http.Request, sub string,
) {
	f.InResponseWriter, f.InReq, f.InSub = w, r, sub
}
