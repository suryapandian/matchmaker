package games

import (
	"errors"
	"sync"

	"github.com/suryapandian/matchmaker/users"
)

type Games struct {
	Store   *sync.Map
	Players users.PlayerStore
}

var gamesStore sync.Map

func NewGame() *Games {
	return &Games{
		Store:   &gamesStore,
		Players: users.NewPlayer(),
	}
}

func (g *Games) save(gameSession *GameSession) {
	g.Store.Store(gameSession.ID, gameSession)
}

func (g *Games) delete(id string) {
	g.Store.Delete(id)

}

var ErrGameNotFound = errors.New("game not found.")

func (g *Games) GetGameByID(id string) (*GameSession, error) {
	game, ok := g.Store.Load(id)
	if !ok {
		return nil, ErrGameNotFound
	}
	return game.(*GameSession), nil
}

func (g *Games) ActiveGames() (games []*GameSession) {
	g.Store.Range(func(_, game interface{}) bool {
		games = append(games, game.(*GameSession))
		return true
	})

	return games
}
