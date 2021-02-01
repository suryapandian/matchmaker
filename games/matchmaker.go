package games

import (
	"errors"
	"github.com/suryapandian/matchmaker/users"

	"github.com/sirupsen/logrus"
)

type Matchmaker struct {
	waitingPlayers  chan string
	stopMatchMaking chan bool
	maximumPlayers  int
	Games           GameStore
	Players         users.PlayerStore
	logger          *logrus.Entry
}

func NewMatchmaker(maximumPlayers int, logger *logrus.Entry) *Matchmaker {
	return &Matchmaker{
		waitingPlayers:  make(chan string, maximumPlayers),
		stopMatchMaking: make(chan bool),
		maximumPlayers:  maximumPlayers,
		Games:           NewGame(),
		Players:         users.NewPlayer(),
		logger:          logger,
	}
}

func (m *Matchmaker) Start() {
	go m.start()
}

func (m *Matchmaker) start() {
	var players []string
	for {
		select {
		case player := <-m.waitingPlayers:
			players = append(players, player)
			m.logger.Infof("Players trying to start %v", players)

			gameSession := NewGameSession(players, "")
			if err := gameSession.start(); err != nil {
				if err == ErrInadequatePlayers {
					continue
				}
				m.logger.WithError(err).Errorln("error starting game session")
			}
			m.logger.WithField("game session id", gameSession.ID).WithField("players", gameSession.Players).Infoln("game session")
			players = []string{}
		case <-m.stopMatchMaking:
			return
		}
	}
}

var (
	ErrMaximumPlayers = errors.New("maximum players are playing in this instance.")
)

func (m *Matchmaker) Join(playerID string) (gameID string, err error) {
	m.logger.WithField("playerId", playerID).Infoln("player trying to join")

	player, err := m.Players.GetPlayerDetails(playerID)
	switch err {
	case users.ErrPlayerNotFound:
		m.Players.ReRegister(playerID)
	case users.ErrNoGame:
		m.logger.WithField("playerId", playerID).Infof("player does not have any on going game")
	case nil:
		m.logger.WithField("playerId", playerID).WithField("gameID", player.GameID).Infoln("player already in game")
		return player.GameID, nil
	default:
		return "", err

	}

	if m.Players.CountActivePlayers() >= m.maximumPlayers {
		return "", ErrMaximumPlayers
	}

	m.logger.WithField("playerId", playerID).Infoln("player waiting to join")
	m.waitingPlayers <- playerID
	return "", ErrInadequatePlayers
}

func (m *Matchmaker) Leave(playerID string) error {
	m.logger.WithField("playerId", playerID).Infoln("player trying to leave")

	player, err := m.Players.GetPlayerDetails(playerID)
	switch err {
	case users.ErrPlayerNotFound:
		return nil
	case users.ErrNoGame:
		m.Players.Deactivate(playerID)
		return nil
	case nil:
	default:
		return err
	}

	if err := m.endGame(player.GameID, playerID); err != nil {
		return err
	}

	return nil
}

func (m *Matchmaker) endGame(id, leavingPlayerID string) error {
	gameSession, err := m.Games.GetGameByID(id)
	if err != nil {
		return err
	}
	m.logger.WithField("gameID", gameSession.ID).Infoln("ending game as player has left")
	gameSession.end()

	for _, player := range gameSession.Players {
		m.logger.WithField("playerId", player).WithField("gameID", gameSession.ID).Infoln("player leaving game")
		if player == leavingPlayerID {
			m.Players.Deactivate(leavingPlayerID)
			continue
		}

		_, err = m.Join(player)
		if err == ErrInadequatePlayers {
			continue
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Matchmaker) Stop() {
	m.stopMatchMaking <- true
}
