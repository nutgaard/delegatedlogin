package server

import (
	. "frontend-image/internal/config"
	. "frontend-image/internal/logging"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type Server struct {
	*fiber.App
}

func CreateServer(config *AppConfig) *Server {
	server := fiber.New(fiber.Config{AppName: config.AppName})
	server.Use(logger.New(logger.Config{
		Output: CreateMaskingWriter(
			"\\d{7,}",
			"*",
			os.Stdout,
		),
	}))

	setupGracefulShutdown(server)

	return &Server{server}
}

func (server Server) Start(port string) {
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}
	if server.Listen(port) != nil {
		panic("Could not start server")
	}
}

func setupGracefulShutdown(server *fiber.App) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig

		log.Info().Msg("Gracefully shutting down...")

		err := server.Shutdown()
		if err != nil {
			log.Fatal().Err(err).Msg("Error shutting down server.")
		}
	}()
}
