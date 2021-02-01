package games

import (
	"testing"

	"github.com/suryapandian/matchmaker/config"

	"github.com/stretchr/testify/suite"
)

type GameTestSuite struct {
	suite.Suite
	games               *Games
	allowedPlayersCount int
}

func (t *GameTestSuite) SetupTest() {
	t.allowedPlayersCount = 3
	config.GAME_CONFIGS[config.GameTypeDefault] = config.GameConfig{AllowedPlayers: t.allowedPlayersCount}
	t.games = NewGame()
}

func (t *GameTestSuite) TestGetActiveGames() {
	game1 := createTestGame(t.allowedPlayersCount)
	game2 := createTestGame(t.allowedPlayersCount)

	games := t.games.ActiveGames()
	var game1Presence, game2Presence bool

	for _, game := range games {
		if !game1Presence {
			game1Presence = (game.ID == game1.ID)
		}

		if !game2Presence {
			game2Presence = game.ID == game2.ID
		}
	}
	t.True(game1Presence)
	t.True(game2Presence)
}

func (t *GameTestSuite) TestGetGame() {
	actualGame := createTestGame(t.allowedPlayersCount)

	expectedGame, err := t.games.GetGameByID(actualGame.ID)
	t.NoError(err)
	t.Equal(actualGame.Players, expectedGame.Players)
	t.Equal(actualGame.ID, expectedGame.ID)
}

func createTestGame(allowedPlayersCount int) *GameSession {
	var players []string
	config.GAME_CONFIGS["default"] = config.GameConfig{AllowedPlayers: allowedPlayersCount}
	games := NewGame()

	for i := 1; i <= allowedPlayersCount; i++ {
		playerID := games.Players.Register()
		players = append(players, playerID)
	}

	game := NewGameSession(players, "")
	game.start()
	return game

}

func TestGameSuite(t *testing.T) {
	suite.Run(t, new(GameTestSuite))
}
