package users

import (
	"errors"
	"sync"

	"github.com/suryapandian/matchmaker/uuid"
)

type Players struct {
	store *sync.Map
}

var store sync.Map

func NewPlayer() *Players {
	return &Players{store: &store}
}

func (p *Players) Register() string {
	id := uuid.NewUUID()
	p.store.Store(id, "")
	return id
}

func (p *Players) isValidID(id string) bool {
	return uuid.IsValidUUID(id)
}

var ErrInvalidPlayerID = errors.New("invalid player id")

func (p *Players) ReRegister(id string) error {
	if !p.isValidID(id) {
		return ErrInvalidPlayerID
	}

	p.store.Store(id, "")
	return nil
}

func (p *Players) Deactivate(id string) {
	p.store.Delete(id)
}

func (p *Players) JoinGame(playerIDs []string, gameID string) {
	for _, id := range playerIDs {
		p.store.Store(id, gameID)
	}
}

func (p *Players) LeaveGame(playerIDs []string) {
	for _, id := range playerIDs {
		p.store.Store(id, "")
	}
}

type Player struct {
	ID     string
	GameID string
}

var (
	ErrPlayerNotFound = errors.New("player not found.")
	ErrNoGame         = errors.New("waiting for other to join game.")
)

func (p *Players) GetPlayerDetails(id string) (*Player, error) {
	if !p.isValidID(id) {
		return nil, ErrInvalidPlayerID
	}

	game, ok := p.store.Load(id)
	if !ok {
		return nil, ErrPlayerNotFound
	}

	gameID := game.(string)
	if gameID == "" {
		return &Player{ID: id}, ErrNoGame
	}

	return &Player{ID: id, GameID: gameID}, nil
}

func (p *Players) GetAvailablePlayers(playerIDs []string) (availablePlayers []string) {
	for _, playerID := range playerIDs {
		_, err := p.GetPlayerDetails(playerID)
		if err == ErrNoGame {
			availablePlayers = append(availablePlayers, playerID)
		}
	}
	return availablePlayers
}

func (p *Players) CountActivePlayers() (count int) {
	p.store.Range(func(_, game interface{}) bool {
		if game.(string) != "" {
			count++
		}
		return true
	})

	return count
}

func (p *Players) getAllPlayers() (ids []string) {
	p.store.Range(func(id, _ interface{}) bool {
		ids = append(ids, id.(string))
		return true
	})
	return ids
}
