package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/suryapandian/matchmaker/config"
	"github.com/suryapandian/matchmaker/games"
	"github.com/suryapandian/matchmaker/handlers"
	"github.com/suryapandian/matchmaker/logger"

	"github.com/sirupsen/logrus"
)

func main() {

	matchmaker := games.NewMatchmaker(
		config.MAX_PLAYERS_PER_INSTANCE,
		logger.LogEntryWithRef(),
	)

	matchmaker.Start()

	logger.SetupLog(config.LOG_LEVEL)
	server := &http.Server{
		Addr:    ":" + config.PORT,
		Handler: handlers.GetRouter(matchmaker),
	}

	go func(server *http.Server) {
		logrus.Infof("Listening on port %s", config.PORT)
		if err := server.ListenAndServe(); err != nil {
			logrus.WithError(err).Fatal("Failed to start server!")
		}
	}(server)

	stopCh := make(chan os.Signal)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)
	<-stopCh
	matchmaker.Stop()
	logrus.Infof("gracefully stopping matchmaker and shutting down server")

	if err := server.Shutdown(context.Background()); err != nil {
		logrus.WithError(err).Fatal("error shutting server down gracefully")
	}
}
