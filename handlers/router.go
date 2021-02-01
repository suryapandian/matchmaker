package handlers

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/suryapandian/matchmaker/games"
	"github.com/suryapandian/matchmaker/logger"
)

func GetRouter(matchmaker *games.Matchmaker) *chi.Mux {
	mux := chi.NewRouter()

	mux.Use(middleware.Logger)

	setPingRoutes(mux)
	newGameRouter(matchmaker, logger.LogEntryWithRef()).setGameRoutes(mux)

	return mux
}
