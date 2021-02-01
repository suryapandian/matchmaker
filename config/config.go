package config

import (
	"os"
)

var (
	PORT      string
	LOG_LEVEL string
)

type GameConfig struct {
	AllowedPlayers int
}

func init() {
	PORT = os.Getenv("PORT")
	if PORT == "" {
		PORT = "3001"
	}
	LOG_LEVEL = os.Getenv("LOG_LEVEL")
	if LOG_LEVEL == "" {
		LOG_LEVEL = "INFO"
	}

}
