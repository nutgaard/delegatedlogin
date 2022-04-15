package server

import (
	"fmt"
	. "frontend-image/internal/config"
	"github.com/gofiber/fiber/v2"
	"path"
)

func (server *Server) SetupLoginRoutes(config *AppConfig) {
	basepath := path.Clean(fmt.Sprintf("/%s/oauth2", config.AppName))

	server.Get(basepath+"/login", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Alive")
	})
	server.Get(basepath+"/callback", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Ready")
	})
	server.Get(basepath+"/whoami", func(ctx *fiber.Ctx) error {
		return ctx.SendString(fmt.Sprintf("Application: %s\nVersion: %s", config.AppName, config.AppVersion))
	})
}
