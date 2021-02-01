package games

import (
	"testing"

	"github.com/suryapandian/matchmaker/config"
	"github.com/suryapandian/matchmaker/users"

	"github.com/stretchr/testify/suite"
)

type GameSessionTestSuite struct {
	suite.Suite
	allowedPlayersCount int
	playerStore         users.PlayerStore
	gameStore           GameStore
}

func (t *GameSessionTestSuite) SetupTest() {
	t.allowedPlayersCount = 2
	config.GAME_CONFIGS[config.GameTypeDefault] = config.GameConfig{AllowedPlayers: t.allowedPlayersCount}
	t.playerStore = users.NewPlayer()
	t.gameStore = NewGame()
}

func (t *GameSessionTestSuite) TestLoadConfig() {
	config.GAME_CONFIGS["football"] = config.GameConfig{AllowedPlayers: 22}

	gameSession := NewGameSession([]string{}, "")
	gameSession.loadConfig()
	t.Equal(config.GameTypeDefault, gameSession.Type, "should set type as default")

	gameSession = NewGameSession([]string{}, "football")
	gameSession.loadConfig()
	t.Equal("football", gameSession.Type, "should set type as football")

}

func (t *GameSessionTestSuite) TestStart() {
	var players []string
	for i := 1; i <= t.allowedPlayersCount; i++ {
		playerID := t.playerStore.Register()
		players = append(players, playerID)

		gameSession := NewGameSession(players, "")
		err := gameSession.start()
		if i != t.allowedPlayersCount {
			t.Equal(ErrInadequatePlayers, err, "should return error when #players is inadequate")
			continue
		}
		expectedGameSession, err := t.gameStore.GetGameByID(gameSession.ID)
		t.NoError(err)
		t.Equal(gameSession.ID, expectedGameSession.ID)
		t.NoError(err, "should start game session")
	}

	oldPlayer := players[0]
	newPlayerID1 := t.playerStore.Register()
	gameSession := NewGameSession([]string{oldPlayer, newPlayerID1}, "")
	err := gameSession.start()
	t.Equal(ErrInadequatePlayers, err, "should not start game for player already in game")
}

func (t *GameSessionTestSuite) TestValidate() {
	playerID := t.playerStore.Register()
	gameSession := NewGameSession([]string{playerID}, "")
	err := gameSession.validate()
	t.Equal(ErrInadequatePlayers, err)

	playerID2 := t.playerStore.Register()
	gameSession = NewGameSession([]string{playerID, playerID2}, "")
	err = gameSession.validate()
	t.NoError(err)
}

func (t *GameSessionTestSuite) TestEndGame() {
	testGameSession := createTestGame(t.allowedPlayersCount)
	t.NotEmpty(testGameSession)

	testGameSession.end()

	gameSession, err := t.gameStore.GetGameByID(testGameSession.ID)
	t.Equal(ErrGameNotFound, err)
	t.Empty(gameSession)
}

func TestGameSessionSuite(t *testing.T) {
	suite.Run(t, new(GameSessionTestSuite))
}
