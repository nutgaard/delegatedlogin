package config

import (
	. "frontend-image/internal/oidc"
	. "frontend-image/internal/resilience"
	"github.com/rs/zerolog/log"
)

type AppData struct {
	Config      *AppConfig
	OidcClient  *OidcClient
	ProxyConfig []ProxyAppConfig
}

func LoadAndCreateServices() *AppData {
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
	log.Info().Msgf("Loaded config: %s", config)
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
		{
			Prefix: "proxy/open-endpoint",
			Url:    "http://localhost:8089",
		},
		{
			Prefix: "proxy/open-endpoint-no-cookie",
			Url:    "http://localhost:8089",
			RewriteDirectives: []string{
				"SET_HEADER Cookie ''",
			},
		},
		{
			Prefix: "proxy/protected-endpoint",
			Url:    "http://localhost:8089",
		},
		{
			Prefix: "proxy/protected-endpoint-with-cookie-rewrite",
			Url:    "http://localhost:8089",
			RewriteDirectives: []string{
				"SET_HEADER Cookie 'ID_token=$cookie{loginapp_ID_token}'",
				"SET_HEADER Authorization '$cookie{loginapp_ID_token}'",
			},
		},
		{
			Prefix: "env-data",
			RewriteDirectives: []string{
				"RESPOND 200 'APP_NAME: $env{APP_NAME}'",
			},
		},
	}
}
