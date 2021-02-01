package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/suryapandian/matchmaker/config"
	"github.com/suryapandian/matchmaker/games"
	"github.com/suryapandian/matchmaker/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestJoinGame(t *testing.T) {
	var testCases = []struct {
		desc               string
		expectedStatusCode int
	}{
		{
			"join player1",
			http.StatusAccepted,
		},
		{
			"join player2",
			http.StatusAccepted,
		},
		{
			"join player3",
			http.StatusInternalServerError,
		},
	}

	config.MAX_PLAYERS_PER_INSTANCE = 2
	config.ALLOWED_PLAYERS_COUNT = 2
	config.GAME_CONFIGS[config.GameTypeDefault] = config.GameConfig{AllowedPlayers: config.ALLOWED_PLAYERS_COUNT}
	matchmaker := games.NewMatchmaker(
		config.MAX_PLAYERS_PER_INSTANCE,
		logger.LogEntryWithRef(),
	)
	matchmaker.Start()

	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			a := assert.New(t)
			r := httptest.NewRequest(http.MethodPost, "/games/join", nil)
			w := httptest.NewRecorder()
			GetRouter(matchmaker).ServeHTTP(w, r)
			response := w.Result()
			a.Equal(testCase.expectedStatusCode, response.StatusCode)
			time.Sleep(2 * time.Second)
		})
	}

}

type GameTestSuite struct {
	suite.Suite
	matchmaker *games.Matchmaker
}

func (t *GameTestSuite) SetupTest() {
	config.MAX_PLAYERS_PER_INSTANCE = 10
	config.ALLOWED_PLAYERS_COUNT = 2
	t.matchmaker = games.NewMatchmaker(
		config.MAX_PLAYERS_PER_INSTANCE,
		logger.LogEntryWithRef(),
	)
	t.matchmaker.Start()
}

func (t *GameTestSuite) TestCookieJoinGame() {
	r1 := httptest.NewRequest(http.MethodPost, "/games/join", nil)
	w1 := httptest.NewRecorder()
	GetRouter(t.matchmaker).ServeHTTP(w1, r1)
	response1 := w1.Result()
	t.Equal(http.StatusAccepted, response1.StatusCode)

	r2 := httptest.NewRequest(http.MethodPost, "/games/join", nil)
	w2 := httptest.NewRecorder()
	GetRouter(t.matchmaker).ServeHTTP(w2, r2)
	response2 := w2.Result()
	t.Equal(http.StatusAccepted, response2.StatusCode)

	//wait for game to get started
	time.Sleep(2 * time.Second)
	//set cookie to the request from the writer
	r2.Header = http.Header{"Cookie": w2.HeaderMap["Set-Cookie"]}

	// Confirm that cookie is present
	cookie, err := r2.Cookie("playerId")
	t.NoError(err)
	t.NotEmpty(cookie.Value)

	//Confirm that the game has started for player2
	w2 = httptest.NewRecorder()
	GetRouter(t.matchmaker).ServeHTTP(w2, r2)
	response2 = w2.Result()
	t.Equal(http.StatusOK, response2.StatusCode)
}

func (t *GameTestSuite) TestLeaveGame() {
	r1 := httptest.NewRequest(http.MethodPost, "/games/join", nil)
	w1 := httptest.NewRecorder()
	GetRouter(t.matchmaker).ServeHTTP(w1, r1)
	response1 := w1.Result()
	t.Equal(http.StatusAccepted, response1.StatusCode)

	r2 := httptest.NewRequest(http.MethodPost, "/games/join", nil)
	w2 := httptest.NewRecorder()
	GetRouter(t.matchmaker).ServeHTTP(w2, r2)
	response2 := w2.Result()
	t.Equal(http.StatusAccepted, response2.StatusCode)

	//wait for game to get started
	time.Sleep(2 * time.Second)
	//set cookie to the request from the writer
	r2 = httptest.NewRequest(http.MethodPost, "/games/leave", nil)
	r2.Header = http.Header{"Cookie": w2.HeaderMap["Set-Cookie"]}

	w2 = httptest.NewRecorder()
	GetRouter(t.matchmaker).ServeHTTP(w2, r2)
	response2 = w2.Result()
	t.Equal(http.StatusOK, response2.StatusCode)

	//requesting to leave without player info should throw error
	r3 := httptest.NewRequest(http.MethodPost, "/games/leave", nil)
	w3 := httptest.NewRecorder()
	GetRouter(t.matchmaker).ServeHTTP(w3, r3)
	response3 := w3.Result()
	t.Equal(http.StatusBadRequest, response3.StatusCode)
}

func (t *GameTestSuite) TestSessions() {
	r1 := httptest.NewRequest(http.MethodPost, "/games/join", nil)
	w1 := httptest.NewRecorder()
	GetRouter(t.matchmaker).ServeHTTP(w1, r1)
	response1 := w1.Result()
	t.Equal(http.StatusAccepted, response1.StatusCode)

	r2 := httptest.NewRequest(http.MethodPost, "/games/join", nil)
	w2 := httptest.NewRecorder()
	GetRouter(t.matchmaker).ServeHTTP(w2, r2)
	response2 := w2.Result()
	t.Equal(http.StatusAccepted, response2.StatusCode)

	//wait for game to get started
	time.Sleep(2 * time.Second)
	r := httptest.NewRequest(http.MethodGet, "/games/sessions", nil)
	w := httptest.NewRecorder()
	GetRouter(t.matchmaker).ServeHTTP(w, r)
	response := w.Result()
	t.Equal(http.StatusOK, response.StatusCode)

}

func TestGameSuite(t *testing.T) {
	suite.Run(t, new(GameTestSuite))
}
