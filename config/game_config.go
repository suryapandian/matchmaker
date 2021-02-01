package config

import (
	"os"
	"strconv"
)

type GameConfig struct {
	AllowedPlayers int
}

var (
	MAX_PLAYERS_PER_INSTANCE int
	ALLOWED_PLAYERS_COUNT    int
	GAME_CONFIGS             map[string]GameConfig
)

const (
	GameTypeDefault = "default"
)

func init() {
	MAX_PLAYERS_PER_INSTANCE, _ = strconv.Atoi(os.Getenv("MAX_PLAYERS_PER_INSTANCE"))
	if MAX_PLAYERS_PER_INSTANCE == 0 {
		MAX_PLAYERS_PER_INSTANCE = 100
	}

	ALLOWED_PLAYERS_COUNT, _ = strconv.Atoi(os.Getenv("ALLOWED_PLAYERS_COUNT"))
	if ALLOWED_PLAYERS_COUNT == 0 {
		ALLOWED_PLAYERS_COUNT = 4
	}

	GAME_CONFIGS = make(map[string]GameConfig)
	GAME_CONFIGS[GameTypeDefault] = GameConfig{AllowedPlayers: ALLOWED_PLAYERS_COUNT}
}
