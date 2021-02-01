package handlers

import (
	"github.com/go-chi/chi"
	"net/http"
)

func setPingRoutes(router chi.Router) {
	router.Route("/", func(r chi.Router) {
		r.Get("/ping", ping)
	})
}

func ping(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("pong!"))
}
