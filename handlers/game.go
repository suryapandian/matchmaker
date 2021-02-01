package handlers

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"net/http"

	"github.com/suryapandian/matchmaker/games"

	"github.com/sirupsen/logrus"
)

type gameRouter struct {
	matchmaker *games.Matchmaker
	logger     *logrus.Entry
}

func newGameRouter(matchmaker *games.Matchmaker, logger *logrus.Entry) *gameRouter {
	return &gameRouter{
		matchmaker: matchmaker,
		logger:     logger,
	}

}

func (g *gameRouter) setGameRoutes(router chi.Router) {
	router.Route("/games", func(r chi.Router) {
		r.Post("/join", g.join)
		r.Post("/leave", g.leave)
		r.Get("/sessions", g.sessions)
	})
}

func (g *gameRouter) join(w http.ResponseWriter, r *http.Request) {
	playerID, _ := g.getPlayerStatusID(r)
	if playerID == "" {
		playerID = g.matchmaker.Players.Register()
		g.logger.Infof("registered the player %s and setting cookie", playerID)
		cookie := http.Cookie{Name: "playerId", Value: playerID, Path: "/", MaxAge: 90000}
		http.SetCookie(w, &cookie)
	}

	gameID, err := g.matchmaker.Join(playerID)

	switch err {
	case games.ErrMaximumPlayers:
		writeJSONMessage(fmt.Sprintf("Please wait, %v", err.Error()), http.StatusInternalServerError, w)
	case games.ErrInadequatePlayers:
		response := fmt.Sprintf("Successfully registered with id: %s. Inadequate number of players to start the game. Please wait!", playerID)
		writeJSONMessage(response, http.StatusAccepted, w)
	case nil:
		response := fmt.Sprintf("Successfully registered with id: %s and joined game with id: %s", playerID, gameID)
		writeJSONMessage(response, http.StatusOK, w)
	default:
		writeJSONMessage(err.Error(), http.StatusInternalServerError, w)

	}

}

func (g *gameRouter) leave(w http.ResponseWriter, r *http.Request) {
	playerID, err := g.getPlayerStatusID(r)
	if playerID == "" {
		writeJSONMessage(err.Error(), http.StatusBadRequest, w)
		return
	}

	err = g.matchmaker.Leave(playerID)
	if err != nil {
		writeJSONMessage(err.Error(), http.StatusInternalServerError, w)
		return
	}

	cookie := http.Cookie{Name: "playerId", Path: "/", MaxAge: -1}
	http.SetCookie(w, &cookie)
	writeJSONMessage("Player has successfully left", http.StatusOK, w)
}

// type req struct {
// 	PlayerID string `json:"playerId"`
// }

func (g *gameRouter) getPlayerStatusID(r *http.Request) (playerID string, err error) {
	ErrPlayNotInGame := errors.New("Player is not playing any game!")
	playerCookie, err := r.Cookie("playerId")
	if err != nil {
		g.logger.WithError(err).Errorln("error reading cookie")
		return "", ErrPlayNotInGame
	}

	playerID = playerCookie.Value
	g.logger.Infof("reading player id from cookie: %s", playerID)
	if playerID == "" {
		return "", ErrPlayNotInGame
	}

	return playerID, nil

	// var req req
	// err = json.NewDecoder(r.Body).Decode(&req)
	// if err != nil {
	// 	return "", err
	// }
	// return req.PlayerID, nil
}

func (g *gameRouter) sessions(w http.ResponseWriter, r *http.Request) {
	games := g.matchmaker.Games.ActiveGames()
	writeJSONStruct(games, http.StatusOK, w)
}
