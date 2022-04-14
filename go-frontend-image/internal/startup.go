package internal

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type AppData struct {
	Config      *AppConfig
	OidcClient  *OidcClient
	ProxyConfig []ProxyAppConfig
}

func LoadAndCreateServices() *AppData {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	config := getConfig()
	oidcClient := getOidcClient(config)

	return &AppData{
		Config:      config,
		OidcClient:  oidcClient,
		ProxyConfig: getProxyConfig(),
	}
}

func getConfig() *AppConfig {
	config := ReadConfig()
	log.Info().Msgf("Loaded config: %v", config)
	return config
}

func getOidcClient(config *AppConfig) *OidcClient {
	var err error
	var oidcClient *OidcClient
	if config.WithoutSecurity {
		return nil
	}

	err = Retry(config.IdpRetryCount, config.IdpRetryDelay, func() error {
		oidcClient, err = CreateOidcClient(
			config.IdpDiscoveryUrl,
			config.IdpClientId,
			config.IdpClientSecret,
		)
		if err != nil {
			log.Warn().Msgf("Could not connect to IDP retrying in %v.\nError: %s", config.IdpRetryDelay, err)
		}
		return err
	})

	if err != nil {
		log.Fatal().Msgf("Could not connect to IDP after %d attempts.\nError: %s", config.IdpRetryCount, err)
	}

	return oidcClient
}

func getProxyConfig() []ProxyAppConfig {
	return []ProxyAppConfig{
		{
			Prefix: "api",
			Url:    "http://localhost:8089/modiapersonoversikt-api",
		},
		{
			Prefix: "proxy/app1",
			Url:    "http://localhost:8089/appname1",
		},
		{
			Prefix: "proxy/app2",
			Url:    "http://localhost:8089/appname2",
		},
	}
}
