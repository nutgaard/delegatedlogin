package config

import (
	"fmt"
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
		ProxyConfig: getProxyConfig(config.DockerCompose),
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

func getProxyConfig(dockerCompose bool) []ProxyAppConfig {
	domain := "localhost:8089"
	if dockerCompose {
		domain = "echo-server"
	}
	return []ProxyAppConfig{
		{
			Prefix: "api",
			Url:    fmt.Sprintf("http://%s/modiapersonoversikt-api", domain),
		},
		{
			Prefix: "proxy/app1",
			Url:    fmt.Sprintf("http://%s/appname1", domain),
		},
		{
			Prefix: "proxy/app2",
			Url:    fmt.Sprintf("http://%s/appname2", domain),
		},
		{
			Prefix: "proxy/open-endpoint",
			Url:    fmt.Sprintf("http://%s", domain),
		},
		{
			Prefix: "proxy/open-endpoint-no-cookie",
			Url:    fmt.Sprintf("http://%s", domain),
			RewriteDirectives: []string{
				"SET_HEADER Cookie ''",
			},
		},
		{
			Prefix: "proxy/protected-endpoint",
			Url:    fmt.Sprintf("http://%s", domain),
		},
		{
			Prefix: "proxy/protected-endpoint-with-cookie-rewrite",
			Url:    fmt.Sprintf("http://%s", domain),
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
