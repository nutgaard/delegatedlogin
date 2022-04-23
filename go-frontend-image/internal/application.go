package internal

import (
	"fmt"
	. "frontend-image/internal/config"
	"frontend-image/internal/router"
	"github.com/rs/zerolog/log"
)

func StartApplication() {
	appData := LoadAndCreateServices()
	config := appData.Config
	oidcClient := appData.OidcClient
	if false {
		fmt.Print(oidcClient) // Just to make it happy
	}

	log.Printf("Starting application: %s (%s)", config.AppName, config.AppVersion)

	r := router.New(appData)
	router.StartServer(config.Port, r)
}
