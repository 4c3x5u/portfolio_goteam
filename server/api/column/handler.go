package column

import (
	"encoding/json"
	"net/http"

	"server/api"
	"server/auth"
	pkgLog "server/log"
)

// Handler is a http.Handler that can be used to handle column requests.
type Handler struct {
	authHeaderReader   auth.HeaderReader
	authTokenValidator auth.TokenValidator
	idValidator        api.StringValidator
	log                pkgLog.Errorer
}

// NewHandler creates and returns a new Handler.
func NewHandler(
	authHeaderReader auth.HeaderReader,
	authTokenValidator auth.TokenValidator,
	idValidator api.StringValidator,
	log pkgLog.Errorer,
) Handler {
	return Handler{
		authHeaderReader:   authHeaderReader,
		authTokenValidator: authTokenValidator,
		idValidator:        idValidator,
		log:                log,
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

	// Get and validate the column ID.
	columnID := r.URL.Query().Get("id")
	if err := h.idValidator.Validate(columnID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err = json.NewEncoder(w).Encode(
			ResBody{Error: err.Error()},
		); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	}

	// All went well. Return 200.
	w.WriteHeader(http.StatusOK)
}
