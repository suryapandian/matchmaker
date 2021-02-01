package users

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type PlayerTestSuite struct {
	suite.Suite
	Players             *Players
	allowedPlayersCount int
}

func (t *PlayerTestSuite) SetupTest() {
	t.Players = NewPlayer()
}

func (t *PlayerTestSuite) TestRegister() {
	playerID := t.Players.Register()
	t.NotEmpty(playerID)
}

func (t *PlayerTestSuite) TestIsValidPlayerID() {
	var testCases = []struct {
		playerID string
		isValid  bool
	}{
		{
			"",
			false,
		},
		{
			"invalid",
			false,
		},
		{
			"497a2922-0d30-41c6-9bd9-c3f8d077de3a",
			true,
		},
	}

	for _, testCase := range testCases {
		validity := t.Players.isValidID(testCase.playerID)
		t.Equal(testCase.isValid, validity)
	}
}

func (t *PlayerTestSuite) TestReRegister() {
	var testCases = []struct {
		playerID string
		err      error
	}{
		{
			"",
			ErrInvalidPlayerID,
		},
		{
			"invalid",
			ErrInvalidPlayerID,
		},
		{
			"497a2922-0d30-41c6-9bd9-c3f8d077de3a",
			nil,
		},
	}

	for _, testCase := range testCases {
		err := t.Players.ReRegister(testCase.playerID)
		t.Equal(testCase.err, err)
	}
}

func (t *PlayerTestSuite) TestDeactivate() {
	playerID := t.Players.Register()
	t.NotEmpty(playerID)

	t.Players.Deactivate(playerID)

	_, err := t.Players.GetPlayerDetails(playerID)
	t.Equal(ErrPlayerNotFound, err)
}

func (t *PlayerTestSuite) TestGetPlayerStatus() {
	_, err := t.Players.GetPlayerDetails("invalid-ID")
	t.Equal(ErrInvalidPlayerID, err)

	_, err = t.Players.GetPlayerDetails("497a2922-0d30-41c6-9bd9-c3f8d077de33")
	t.Equal(ErrPlayerNotFound, err)

	playerID := t.Players.Register()
	t.NotEmpty(playerID)

	_, err = t.Players.GetPlayerDetails(playerID)
	t.Equal(ErrNoGame, err)
}

func (t *PlayerTestSuite) TestJoinGame() {
	playerID := t.Players.Register()
	t.NotEmpty(playerID)

	gameID := "497a2922-0d30-41c6-9bd9-c3f8d077de3a"
	t.Players.JoinGame([]string{playerID}, gameID)

	player, err := t.Players.GetPlayerDetails(playerID)
	t.NoError(err)
	t.Equal(gameID, player.GameID)
}

func (t *PlayerTestSuite) TestLeaveGame() {
	playerID := t.Players.Register()
	t.NotEmpty(playerID)

	gameID := "497a2922-0d30-41c6-9bd9-c3f8d077de3a"
	t.Players.JoinGame([]string{playerID}, gameID)

	player, err := t.Players.GetPlayerDetails(playerID)
	t.Equal(gameID, player.GameID)

	t.Players.LeaveGame([]string{playerID})

	player, err = t.Players.GetPlayerDetails(playerID)
	t.Equal(ErrNoGame, err)
	t.Empty(player.GameID)
}

func (t *PlayerTestSuite) TestGetAvailablePlayers() {
	playerID := t.Players.Register()
	t.NotEmpty(playerID)

	gameID := "497a2922-0d30-41c6-9bd9-c3f8d077de3a"
	t.Players.JoinGame([]string{playerID}, gameID)

	playerID2 := t.Players.Register()
	t.NotEmpty(playerID2)

	availablePlayers := t.Players.GetAvailablePlayers([]string{playerID, playerID2})
	t.Equal([]string{playerID2}, availablePlayers)
}

func (t *PlayerTestSuite) TestCountActivePlayers() {
	//cleanup
	playerIDs := t.Players.getAllPlayers()
	for _, id := range playerIDs {
		t.Players.Deactivate(id)
	}

	playerID := t.Players.Register()
	t.NotEmpty(playerID)

	playerID2 := t.Players.Register()
	t.NotEmpty(playerID2)

	totalPlayers := t.Players.CountActivePlayers()
	t.Equal(0, totalPlayers, "registered players not playing any games")

	gameID := "497a2922-0d30-41c6-9bd9-c3f8d077de3a"
	t.Players.JoinGame([]string{playerID, playerID2}, gameID)

	totalPlayers = t.Players.CountActivePlayers()
	t.Equal(2, totalPlayers, "players actively playing any games")
}

func TestPlayerSuite(t *testing.T) {
	suite.Run(t, new(PlayerTestSuite))
}
