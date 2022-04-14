package internal

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"strings"
)

type ProxyAppConfig struct {
	Prefix string
	Url    string
}

func (server *Server) SetupProxyRoutes(prefix string, config ...ProxyAppConfig) {
	keySelector := func(app ProxyAppConfig) string { return app.Prefix }
	trie := NewPathTrie(config, keySelector)

	handler := func(ctx *fiber.Ctx) error {
		url := strings.TrimPrefix(ctx.Path(), prefix)
		urlWithoutLeadingSlash := strings.TrimPrefix(url, "/")
		appconfig := trie.Search(urlWithoutLeadingSlash)

		if appconfig != nil {
			proxyUrl := appconfig.Url + strings.TrimPrefix(urlWithoutLeadingSlash, appconfig.Prefix)
			ctx.Request().Header.Set("Authorization", "Bearer dummy-token")
			return proxy.Do(ctx, proxyUrl)
		} else {
			return ctx.Next()
		}
	}

	server.Use(handler)
}
