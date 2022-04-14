package internal

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"path"
)

type jsonData map[string]interface{}

func (server *Server) SetupApiRoutes(config *AppConfig) {
	apipath := path.Clean(fmt.Sprintf("/%s/self", config.AppName))
	server.Get(apipath+"/hello", func(ctx *fiber.Ctx) error {
		return ctx.JSON(
			jsonData{
				"id":   "1234",
				"name": "Name Nameson",
			},
		)
	})

}
