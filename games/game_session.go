package games

import (
	"errors"

	"github.com/suryapandian/matchmaker/array"
	"github.com/suryapandian/matchmaker/config"
	"github.com/suryapandian/matchmaker/users"

	"github.com/suryapandian/matchmaker/uuid"
)

type GameSession struct {
	ID      string            `json:"id"`
	Players []string          `json:"players"`
	Type    string            `json:"-"`
	Config  config.GameConfig `json:"-"`
}

func NewGameSession(players []string, gameType string) *GameSession {
	players = array.RemoveDuplicates(players)
	playerStore := users.NewPlayer()
	gameSession := &GameSession{
		Players: playerStore.GetAvailablePlayers(players),
		Type:    gameType,
	}
	return gameSession
}

func (gameSession *GameSession) start() error {
	if err := gameSession.validate(); err != nil {
		return err
	}
	gameSession.ID = uuid.NewUUID()

	NewGame().save(gameSession)
	users.NewPlayer().JoinGame(gameSession.Players, gameSession.ID)
	return nil
}

var ErrInadequatePlayers = errors.New("inadequate players to start a game session")

func (gameSession *GameSession) validate() error {
	gameSession.loadConfig()
	if len(gameSession.Players) < gameSession.Config.AllowedPlayers {
		return ErrInadequatePlayers
	}
	return nil
}

func (gameSession *GameSession) loadConfig() {
	if gameSession.Type == "" {
		gameSession.Type = config.GameTypeDefault
	}

	gameSession.Config = config.GAME_CONFIGS[gameSession.Type]
}

func (gameSession *GameSession) end() {
	NewGame().delete(gameSession.ID)
	users.NewPlayer().LeaveGame(gameSession.Players)
}
