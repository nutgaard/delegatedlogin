package oidc

import (
	"context"
	"encoding/json"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jws"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/url"
	"time"
)

type OidcClient struct {
	discoveryUrl string
	clientId     string
	clientSecret string
	JwksConfig   *JwksConfig
	keySet       jwk.Set
}
type JwksConfig struct {
	JwksUri               string `json:"jwks_uri"`
	TokenEndpoint         string `json:"token_endpoint"`
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	Issuer                string `json:"issuer"`
}
type TokenExchangeResult struct {
	IdToken      string `json:"id_token"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type RefreshIdTokenResponse struct {
	IdToken string `json:"id_token"`
}
type RefreshIdTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

var httpClient = http.Client{}

func (client OidcClient) OpenAmExchangeAuthCodeForToken(code string, loginurl string) (*TokenExchangeResult, error) {
	req, err := http.NewRequest(
		"POST",
		client.JwksConfig.TokenEndpoint,
		nil,
	)
	req.SetBasicAuth(client.clientId, client.clientSecret)
	req.PostForm = url.Values{}
	req.PostForm.Add("grant_type", "authorization_code")
	req.PostForm.Add("realm", "/")
	req.PostForm.Add("redirect_uri", loginurl)
	req.PostForm.Add("code", code)
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	var result *TokenExchangeResult
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (client OidcClient) RefreshIdToken(refreshToken string) (*RefreshIdTokenResponse, error) {
	req, err := http.NewRequest(
		"POST",
		client.JwksConfig.TokenEndpoint,
		nil,
	)
	req.SetBasicAuth(client.clientId, client.clientSecret)
	req.PostForm = url.Values{}
	req.PostForm.Add("grant_type", "refresh_token")
	req.PostForm.Add("realm", "/")
	req.PostForm.Add("scope", "openid")
	req.PostForm.Add("refresh_token", refreshToken)
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	var result *RefreshIdTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (client OidcClient) Verify(token string) (jwt.Token, error) {
	return jwt.Parse(
		[]byte(token),
		jwt.WithKeySet(client.keySet, jws.WithInferAlgorithmFromKey(true)),
		jwt.WithAudience(client.clientId),
		jwt.WithVerify(true),
		jwt.WithValidate(true),
	)
}

func CreateOidcClient(discoveryUrl, clientId, clientSecret string) (*OidcClient, error) {
	jwksConfig, err := fetchConfig(discoveryUrl)
	if err != nil {
		return nil, err
	}
	keySet, err := fetchKeys(jwksConfig.JwksUri)
	return &OidcClient{
		discoveryUrl: discoveryUrl,
		clientId:     clientId,
		clientSecret: clientSecret,
		JwksConfig:   jwksConfig,
		keySet:       keySet,
	}, nil
}

func fetchConfig(discoveryUrl string) (*JwksConfig, error) {
	req, err := http.NewRequest("GET", discoveryUrl, nil)
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Accept-Type", "application/json")

	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	var config *JwksConfig
	err = json.NewDecoder(resp.Body).Decode(&config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func fetchKeys(keysUrl string) (jwk.Set, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cache := jwk.NewCache(ctx)
	err := cache.Register(keysUrl, jwk.WithMinRefreshInterval(15*time.Minute))
	if err != nil {
		return nil, err
	}

	_, err = cache.Refresh(ctx, keysUrl)
	if err != nil {
		log.Err(err).Msg("failed to refresh jwks")
		return nil, err
	}
	set := jwk.NewCachedSet(cache, keysUrl)
	return set, nil
}
