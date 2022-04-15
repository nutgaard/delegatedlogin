package server

import (
	"github.com/gofiber/fiber/v2"
	"path"
)

func (server *Server) SetupStaticFileRoutes(prefix, root string) {
	indexFile := path.Join(root, "index.html")
	server.Static(prefix, root)
	server.Get("*", func(ctx *fiber.Ctx) error {
		return ctx.SendFile(indexFile)
	})
}
