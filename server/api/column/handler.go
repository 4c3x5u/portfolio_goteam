package column

import (
	"net/http"
	"server/api"

	"server/auth"
)

// Handler is a http.Handler that can be used to handle column requests.
type Handler struct {
	authHeaderReader   auth.HeaderReader
	authTokenValidator auth.TokenValidator
}

// NewHandler creates and returns a new Handler.
func NewHandler(
	authHeaderReader auth.HeaderReader,
	authTokenValidator auth.TokenValidator,
) Handler {
	return Handler{
		authHeaderReader:   authHeaderReader,
		authTokenValidator: authTokenValidator,
	}
}

// ServeHTTP responds to requests made to the column route.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only allow PATCH requests.
	if r.Method != http.MethodPatch {
		w.Header().Add(api.AllowedMethods(http.MethodPost))
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	// Get auth token from Authorization header, validate it, and get
	// the subject of the token.
	authToken := h.authHeaderReader.Read(
		r.Header.Get(auth.AuthorizationHeader),
	)
	sub := h.authTokenValidator.Validate(authToken)
	if sub == "" {
		w.Header().Set(auth.WWWAuthenticate())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// All went well. Return 200.
	w.WriteHeader(http.StatusOK)
}
