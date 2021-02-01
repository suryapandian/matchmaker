package games

import (
	"testing"
	"time"

	"github.com/suryapandian/matchmaker/config"
	"github.com/suryapandian/matchmaker/logger"
	"github.com/suryapandian/matchmaker/users"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MatchMakerTestSuite struct {
	suite.Suite
	matchmaker          *Matchmaker
	maximumPlayers      int
	allowedPlayersCount int
}

func (t *MatchMakerTestSuite) SetupTest() {
	t.maximumPlayers = 9
	t.allowedPlayersCount = 3
	config.GAME_CONFIGS[config.GameTypeDefault] = config.GameConfig{AllowedPlayers: t.allowedPlayersCount}
	t.matchmaker = NewMatchmaker(
		t.maximumPlayers,
		logger.LogEntryWithRef(),
	)
	t.matchmaker.Start()
}

func (t *MatchMakerTestSuite) TestJoin() {
	clearGames(t.matchmaker)

	_, err := t.matchmaker.Join("", "")
	t.Equal(users.ErrInvalidPlayerID, err, "should register and join as new player")

	var players []string

	for i := 0; i < t.maximumPlayers; i++ {
		playerID := t.matchmaker.Players.Register()
		players = append(players, playerID)

		t.matchmaker.Join(playerID, "")
	}

	time.Sleep(2 * time.Second)
	for i := 0; i < t.maximumPlayers; i++ {
		gameID, err := t.matchmaker.Join(players[i], "")
		t.NoError(err, "should start game")
		t.NotEmpty(gameID)
	}

	playerID := t.matchmaker.Players.Register()
	_, err = t.matchmaker.Join(playerID, "")
	t.Equal(ErrMaximumPlayers, err, "should throw error on reaching max players per instance")
}

func (t *MatchMakerTestSuite) TestLeave() {
	var players []string
	clearGames(t.matchmaker)

	for i := 0; i < t.allowedPlayersCount; i++ {
		playerID := t.matchmaker.Players.Register()
		players = append(players, playerID)

		t.matchmaker.Join(playerID, "")
	}

	time.Sleep(2 * time.Second)
	gameID, err := t.matchmaker.Join(players[0], "")
	t.NoError(err, "should start game session")

	game, err := t.matchmaker.Games.GetGameByID(gameID)
	t.NoError(err, "should get game session details")
	t.NotEmpty(game)

	err = t.matchmaker.Leave(players[0])
	t.NoError(err, "should leave game")

	player, err := t.matchmaker.Players.GetPlayerDetails(players[0])
	t.Equal(err, users.ErrPlayerNotFound, "player who has left the game should not be found")
	t.Nil(player)

	for i := 1; i < len(players); i++ {
		player, err = t.matchmaker.Players.GetPlayerDetails(players[1])
		t.Equal(err, users.ErrNoGame, "other player should be in queue without any active session")
	}

	game, err = t.matchmaker.Games.GetGameByID(gameID)
	t.Equal(ErrGameNotFound, err, "game should have quit successfully")
	t.Empty(game)

	//Check if other players are in queue by joining just one another player
	playerID := t.matchmaker.Players.Register()
	t.matchmaker.Join(playerID, "")

	time.Sleep(2 * time.Second)
	gameID, err = t.matchmaker.Join(playerID, "")
	t.NoError(err, "new game should have started for the other players who have not left")
	t.NotEmpty(gameID)
}

func clearGames(matchmaker *Matchmaker) {
	games := matchmaker.Games.ActiveGames()
	for _, game := range games {
		for _, player := range game.Players {
			matchmaker.Leave(player)
		}
	}

}

func (t *MatchMakerTestSuite) TearDownSuite() {
	t.matchmaker.Stop()
}

func TestMatchMakerSuite(t *testing.T) {
	suite.Run(t, new(MatchMakerTestSuite))
}

func TestMultiTypeMatchMaker(t *testing.T) {
	maximumPlayers := 6
	testGameType := "testGameType"
	config.GAME_CONFIGS[config.GameTypeDefault] = config.GameConfig{AllowedPlayers: 3}
	config.GAME_CONFIGS[testGameType] = config.GameConfig{AllowedPlayers: 3}

	matchmaker := NewMatchmaker(
		maximumPlayers,
		logger.LogEntryWithRef(),
	)
	matchmaker.Start()

	clearGames(matchmaker)

	var players []string
	for i := 0; i < config.GAME_CONFIGS[config.GameTypeDefault].AllowedPlayers; i++ {
		playerID := matchmaker.Players.Register()
		players = append(players, playerID)

		matchmaker.Join(playerID, config.GameTypeDefault)
	}

	a := assert.New(t)
	time.Sleep(2 * time.Second)
	_, err := matchmaker.Join(players[0], "")
	a.NoError(err, "should start game session")

	players = []string{}
	for i := 0; i < config.GAME_CONFIGS[testGameType].AllowedPlayers; i++ {
		playerID := matchmaker.Players.Register()
		players = append(players, playerID)

		matchmaker.Join(playerID, testGameType)
	}

	time.Sleep(2 * time.Second)
	_, err = matchmaker.Join(players[0], testGameType)
	a.NoError(err, "should start game session")

	playerID := matchmaker.Players.Register()
	_, err = matchmaker.Join(playerID, config.GameTypeDefault)
	a.Equal(ErrMaximumPlayers, err, "should throw error on reaching max players per instance")

}
