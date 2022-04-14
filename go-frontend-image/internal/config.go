package internal

import (
	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

type AppConfig struct {
	AppName    string `env:"APP_NAME,notEmpty"`
	AppVersion string `env:"APP_VERSION,notEmpty"`

	IdpDiscoveryUrl string        `env:"IDP_DISCOVERY_URL,notEmpty"`
	IdpClientId     string        `env:"IDP_CLIENT_ID,notEmpty"`
	IdpClientSecret string        `env:"IDP_CLIENT_SECRET,notEmpty"`
	IdpRetryCount   int           `env:"IDP_RETRY_COUNT" envDefault:"2"`
	IdpRetryDelay   time.Duration `env:"IDP_RETRY_DELAY_MS" envDefault:"5s"`

	AuthTokenResolver    string `env:"AUTH_TOKEN_RESOLVER,notEmpty"`
	RefreshTokenResolver string `env:"REFRESH_TOKEN_RESOLVER,notEmpty"`

	ReferrerPolicy string `env:"REFERRER_POLICY" envDefault:"origin"`
	CspDirectives  string `env:"CSP_DIRECTIVES" envDefault:"default-src 'self';"`
	CspReportOnly  bool   `env:"CSP_REPORT_ONLY" envDefault:"false"`

	Port string `env:"EXPOSED_PORT" envDefault:"8080"`

	WithoutSecurity bool `env:"WITHOUT_SECURITY" envDefault:"false"`
}

func SetupEnv() {
	_ = os.Setenv("APP_NAME", "modiapersonoversikt")
	_ = os.Setenv("APP_VERSION", "localhost")

	_ = os.Setenv("IDP_DISCOVERY_URL", "http://localhost:8080/.well-known/openid-configuration")
	_ = os.Setenv("IDP_CLIENT_ID", "foo")
	_ = os.Setenv("IDP_CLIENT_SECRET", "bar")

	_ = os.Setenv("AUTH_TOKEN_RESOLVER", "loginapp_ID_token")
	_ = os.Setenv("REFRESH_TOKEN_RESOLVER", "loginapp_refresh_token")

	_ = os.Setenv("WITHOUT_SECURITY", "true")
}

func ReadConfig() *AppConfig {
	config := &AppConfig{}
	if err := env.Parse(config); err != nil {
		log.Fatal().Msg(err.Error())
	}
	return config
}
