package board

import (
	"encoding/json"
	"net/http"

	"server/auth"
	"server/db"
	"server/relay"
)

// Handler is the http.Handler for the boards route.
type Handler struct {
	tokenValidator   auth.TokenValidator
	userBoardCounter db.Counter
	boardInserter    db.Inserter[db.Board]
}

// NewHandler creates and returns a new Handler.
func NewHandler(
	tokenValidator auth.TokenValidator,
	userBoardCounter db.Counter,
	boardInserter db.Inserter[db.Board],
) Handler {
	return Handler{
		tokenValidator:   tokenValidator,
		userBoardCounter: userBoardCounter,
		boardInserter:    boardInserter,
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
	if boardCount := h.userBoardCounter.Count(userID); boardCount >= maxBoards {
		relay.ClientErr(w, http.StatusBadRequest, ResBody{Error: errMaxBoards})
		return
	}

	// Read the request body.
	reqBody := ReqBody{}
	if err = json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		relay.ServerErr(w, err.Error())
		return
	}

	// Create a new board.
	if err := h.boardInserter.Insert(
		db.NewBoard(reqBody.Name, userID),
	); err != nil {
		relay.ServerErr(w, err.Error())
		return
	}
}

// maxBoards is the amount of boards that each user is allowed to own (i.e. be
// the admin of).
const maxBoards = 3

// errMaxBoards is the error message returned from the handler when the user
// already owns the maximum amount of boards allowed per user.
const errMaxBoards = "You have already created the maximum amount of boards " +
	"allowed per user. Please delete one of your boards to create a new one."
