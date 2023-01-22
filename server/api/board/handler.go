package board

import (
	"net/http"

	"server/auth"
	"server/db"
	"server/relay"
)

// Handler is the http.Handler for the boards route.
type Handler struct {
	tokenValidator   auth.TokenValidator
	userBoardCounter db.Counter
}

// NewHandler creates and returns a new Handler.
func NewHandler(
	tokenValidator auth.TokenValidator,
	userBoardCounter db.Counter,
) Handler {
	return Handler{
		tokenValidator:   tokenValidator,
		userBoardCounter: userBoardCounter,
	}
}

// ServeHTTP responds to requests made to the board route.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Accept only POST requests.
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Get the authentication cookie.
	authCookie, err := r.Cookie(auth.CookieName)
	if err != nil {
		w.Header().Set(auth.WWWAuthenticate())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Validate authentication cookie value and get userID.
	userID, err := h.tokenValidator.Validate(authCookie.Value)
	if err != nil {
		w.Header().Set(auth.WWWAuthenticate())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Validate that the user has less than 3 boards. This is done to limit the
	// resources used by this demo app.
	if boardCount := h.userBoardCounter.Count(userID); boardCount >= 3 {
		relay.ClientErr(w, http.StatusBadRequest, ResBody{Error: errMaxBoards})
		return
	}
}

// errMaxBoards is the error message returned from the handler when the user
// already owns the maximum amount of boards allowed per user.
const errMaxBoards = "You have already created the maximum amount of boards " +
	"allowed per user. Please delete one of your boards to create a new one."
