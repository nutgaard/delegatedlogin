package internal

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"os"
	"strings"
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
