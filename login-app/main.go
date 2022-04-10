package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	ConfigLoader "login-app/config"
	"login-app/oidc"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	//ConfigLoader.SetupEnv()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	config := ConfigLoader.ReadConfig()
	log.Info().Msgf("Loaded config: %s", config)

	oidcClient, err := oidc.CreateOIDC(config.IdpDiscoveryUrl, config.IdpClientId, config.IdpClientSecret)
	if err != nil {
		log.Info().Msg("Could not connect to IDP, waiting for 10s for retrying: " + err.Error())
		time.Sleep(10 * time.Second)
	}
	oidcClient, err = oidc.CreateOIDC(config.IdpDiscoveryUrl, config.IdpClientId, config.IdpClientSecret)
	if err != nil {
		log.Fatal().Msg("Could not connect to IDP after 10s: " + err.Error())
	}
	log.Info().Msg("Created OIDC client")

	LoginRoutes(config, oidcClient)
	NaisRoutes(config, oidcClient)

	port := "8080"
	if os.Getenv("DOCKER_COMPOSE") == "false" {
		port = strconv.Itoa(int(config.ExposedPort))
	}
	log.Info().Msg("Listening to " + port)

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
}
