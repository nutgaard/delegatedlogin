package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog/log"
	"os"
)

type Config struct {
	AppName              string `env:"APP_NAME,notEmpty"`
	AppVersion           string `env:"APP_VERSION,notEmpty"`
	IdpDiscoveryUrl      string `env:"IDP_DISCOVERY_URL,notEmpty"`
	IdpClientId          string `env:"IDP_CLIENT_ID,notEmpty"`
	IdpClientSecret      string `env:"IDP_CLIENT_SECRET,notEmpty"`
	AuthTokenResolver    string `env:"AUTH_TOKEN_RESOLVER,notEmpty"`
	RefreshTokenResolver string `env:"REFRESH_TOKEN_RESOLVER,notEmpty"`
	ExposedPort          uint16 `env:"EXPOSED_PORT" envDefault:"8080"`
}

func SetupEnv() {
	_ = os.Setenv("APP_NAME", "loginapp")
	_ = os.Setenv("APP_VERSION", "localhost")
	_ = os.Setenv("IDP_DISCOVERY_URL", "http://localhost:8080/.well-known/openid-configuration")
	_ = os.Setenv("IDP_CLIENT_ID", "foo")
	_ = os.Setenv("IDP_CLIENT_SECRET", "bar")
	_ = os.Setenv("AUTH_TOKEN_RESOLVER", "loginapp_ID_token")
	_ = os.Setenv("REFRESH_TOKEN_RESOLVER", "loginapp_refresh_token")
	_ = os.Setenv("DOCKER_COMPOSE", "false")
	_ = os.Setenv("EXPOSED_PORT", "8082")
}

func ReadConfig() *Config {
	config := &Config{}
	if err := env.Parse(config); err != nil {
		log.Fatal().Msg(err.Error())
	}
	return config
}
