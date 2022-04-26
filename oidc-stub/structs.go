package main

import (
	"github.com/caarlos0/env/v6"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/rs/zerolog/log"
)

type OIDCConfig struct {
	Url                   string `json:"jwks_uri"`
	TokenEndpoint         string `json:"token_endpoint"`
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	Issuer                string `json:"issuer"`
}
type TokenExchangeResult struct {
	IdToken      string `json:"id_token"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AppContext struct {
	Jwks   jwk.Set
	Config *Config
}

type Config struct {
	DockerCompose bool `env:"DOCKER_COMPOSE" envDefault:"false"`
}

func loadConfig() *Config {
	config := &Config{}
	if err := env.Parse(config); err != nil {
		log.Fatal().Err(err).Msg("Could not load config")
	}
	return config
}
