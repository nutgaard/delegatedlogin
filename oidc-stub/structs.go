package main

import "github.com/lestrrat-go/jwx/jwk"

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
	PublicKey  jwk.RSAPublicKey
	PrivateKey jwk.RSAPrivateKey
	Jwks       jwk.Set
	JwksJson   []byte
}
