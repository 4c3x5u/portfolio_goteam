package column

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"server/dbaccess"
	columnTable "server/dbaccess/column"

	"server/api"
	"server/auth"
	pkgLog "server/log"
)

// Handler is a http.Handler that can be used to handle column requests.
type Handler struct {
	authHeaderReader   auth.HeaderReader
	authTokenValidator auth.TokenValidator
	idValidator        api.StringValidator
	columnSelector     dbaccess.Selector[columnTable.Record]
	log                pkgLog.Errorer
}

// NewHandler creates and returns a new Handler.
func NewHandler(
	authHeaderReader auth.HeaderReader,
	authTokenValidator auth.TokenValidator,
	idValidator api.StringValidator,
	columnSelector dbaccess.Selector[columnTable.Record],
	log pkgLog.Errorer,
) Handler {
	return Handler{
		authHeaderReader:   authHeaderReader,
		authTokenValidator: authTokenValidator,
		idValidator:        idValidator,
		columnSelector:     columnSelector,
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

	// Retrieve the column from the database so that we find out its board ID to
	// validate that the user has the right to edit it.
	_, err := h.columnSelector.Select(columnID)
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusBadRequest)
		if err = json.NewEncoder(w).Encode(
			ResBody{Error: "Column not found."},
		); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Error(err.Error())
		}
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(err.Error())
		return
	}

	// All went well. Return 200.
	w.WriteHeader(http.StatusOK)
}
