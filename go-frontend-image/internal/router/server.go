package router

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func StartServer(port string, router chi.Router) {
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	setupGracefulShutdown(server)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal().Err(err).Msgf("Could not start server")
	}
	log.Info().Msgf("Server started")
}

func setupGracefulShutdown(server *http.Server) {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-done
		shutdownCtx, shutdownCancelCtx := context.WithTimeout(context.Background(), 5*time.Second)

		log.Info().Msg("Gracefully shutting down...")

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal().Err(err).Msg("Error shutting down server.")
		}
		shutdownCancelCtx()
	}()
}
