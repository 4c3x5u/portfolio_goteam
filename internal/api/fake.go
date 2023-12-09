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
// the fake instance.
func (f *FakeMethodHandler) Handle(
	w http.ResponseWriter, r *http.Request, sub string,
) {
	f.InResponseWriter, f.InReq, f.InSub = w, r, sub
}

// FakeStringValidator is a test fake for StringValidator.
type FakeStringValidator struct{ Err error }

// Validate implements the StringValidator interface on fakeStringValidator.
func (f *FakeStringValidator) Validate(_ string) error { return f.Err }

// FakeStringValidator is a test fake for IntValidator.
type FakeIntValidator struct{ Err error }

// Validate implements the StringValidator interface on fakeIntValidator.
func (f *FakeIntValidator) Validate(_ int) error { return f.Err }
