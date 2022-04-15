package server

import (
	"fmt"
	. "frontend-image/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"path"
)

func (server *Server) SetupK8sRoutes(config *AppConfig) {
	internalpath := path.Clean(fmt.Sprintf("/%s/internal", config.AppName))

	server.Get(internalpath+"/isAlive", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Alive")
	})
	server.Get(internalpath+"/isReady", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Ready")
	})
	server.Get(internalpath+"/selftest", func(ctx *fiber.Ctx) error {
		return ctx.SendString(fmt.Sprintf("Application: %s\nVersion: %s", config.AppName, config.AppVersion))
	})

	// TODO remove
	server.Get(internalpath+"/dashboard", monitor.New())
}
