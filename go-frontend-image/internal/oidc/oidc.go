package oidc

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type OidcClient struct {
	discoveryUrl string
	clientId     string
	clientSecret string
	JwksConfig   *JwksConfig
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

func CreateOidcClient(discoveryUrl, clientId, clientSecret string) (*OidcClient, error) {
	jwksConfig, err := fetchConfig(discoveryUrl)
	if err != nil {
		return nil, err
	}
	return &OidcClient{
		discoveryUrl: discoveryUrl,
		clientId:     clientId,
		clientSecret: clientSecret,
		JwksConfig:   jwksConfig,
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
