package internal

import (
	"fmt"
	. "frontend-image/internal/config"
	. "frontend-image/internal/middleware"
	. "frontend-image/internal/server"
	"github.com/rs/zerolog/log"
)

func StartApplication() {
	appData := LoadAndCreateServices()
	config := appData.Config
	oidcClient := appData.OidcClient
	proxyConfig := appData.ProxyConfig
	if false {
		fmt.Print(oidcClient) // Just to make it happy
	}

	appPath := fmt.Sprintf("/%s/", config.AppName)
	log.Printf("Starting application: %s (%s)", config.AppName, config.AppVersion)

	server := CreateServer(config)
	server.Use(CallIdMiddleware())
	server.Use(ZerologMiddleware(MaskingConfig{
		Pattern:     "\\d{7,}",
		Replacement: "*",
	}))

	server.SetupK8sRoutes(config)
	server.SetupLoginRoutes(config)
	server.SetupApiRoutes(config)
	server.SetupProxyRoutes(
		appPath,
		proxyConfig...,
	)
	server.SetupStaticFileRoutes(appPath, "./www")

	server.Start(config.Port)
}
