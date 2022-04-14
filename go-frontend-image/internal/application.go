package internal

import (
	"fmt"
)

func StartApplication() {
	appData := LoadAndCreateServices()
	config := appData.Config
	oidcClient := appData.OidcClient
	proxyConfig := appData.ProxyConfig
	appPath := fmt.Sprintf("/%s/", config.AppName)

	fmt.Println(oidcClient)

	server := CreateServer(config)
	server.SetupK8sRoutes(config)
	server.SetupApiRoutes(config)
	server.SetupProxyRoutes(
		appPath,
		proxyConfig...,
	)
	server.SetupStaticFileRoutes(appPath, "./www")

	server.Start(config.Port)
}
