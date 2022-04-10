package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"login-app/accesslog"
	"login-app/config"
	"login-app/oidc"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Context struct {
	config     *config.Config
	oidcClient *oidc.Client
}

func LoginRoutes(config *config.Config, oidcClient *oidc.Client) {
	rand.Seed(time.Now().UnixNano())
	context := Context{config, oidcClient}

	http.HandleFunc(fmt.Sprintf("/%s/api/start", config.AppName), accesslog.Decorate(context.startHandler))
	http.HandleFunc(fmt.Sprintf("/%s/api/login", config.AppName), accesslog.Decorate(context.loginHandler))
	http.HandleFunc(fmt.Sprintf("/%s/api/refresh", config.AppName), accesslog.Decorate(context.refreshHandler))
}

func (context Context) startHandler(w http.ResponseWriter, r *http.Request) {
	returnUrl := r.URL.Query().Get("url")
	if len(returnUrl) == 0 {
		http.Error(w, "URL parameter 'url' is missing", 400)
		return
	}
	stateNounce := createStateNounce()

	setCookie(w, r, stateNounce, returnUrl, 3600)

	callbackUrl, err := createCallbackUrl(
		context.oidcClient.JwksConfig.AuthorizationEndpoint,
		context.config.IdpClientId,
		stateNounce,
		loginUrl(r, context.config.ExposedPort, context.config.AppName),
	)
	if err != nil {
		http.Error(w, "Could not create callback url", 500)
		return
	}

	http.Redirect(w, r, callbackUrl, http.StatusFound)
}

func (context Context) loginHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if len(code) == 0 {
		http.Error(w, "URL parameter 'code' is missing", 400)
		return
	}
	state := r.URL.Query().Get("state")
	if len(state) == 0 {
		http.Error(w, "URL parameter 'state' is missing", 400)
		return
	}
	cookie, err := r.Cookie(state)
	if err != nil {
		http.Error(w, "State-cookie is missing", 400)
		return
	}

	token, err := context.oidcClient.OpenAmExchangeAuthCodeForToken(
		code,
		loginUrl(r, context.config.ExposedPort, context.config.AppName),
	)
	if err != nil {
		http.Error(w, "Could not exchange auth-code for token: "+err.Error(), 500)
	}

	setCookie(w, r, context.config.AuthTokenResolver, token.IdToken, 3600)
	if len(token.RefreshToken) > 0 {
		setCookie(w, r, context.config.RefreshTokenResolver, token.RefreshToken, 20*3600)
	}
	removeCookie(w, r, state)
	originalUrl, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		http.Error(w, "Could not deserialize originalUrl: "+cookie.Value, 500)
	}
	http.Redirect(w, r, originalUrl, http.StatusFound)
}

func (context Context) refreshHandler(w http.ResponseWriter, r *http.Request) {
	var body oidc.RefreshIdTokenRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Could not get refresh token", 500)
		return
	}
	token, err := context.oidcClient.RefreshIdToken(body.RefreshToken)
	if err != nil {
		http.Error(w, "Could not refresh token: "+body.RefreshToken, 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(oidc.RefreshIdTokenResponse{
		IdToken: token.IdToken,
	})
	if err != nil {
		http.Error(w, "Could not serialize refresh token"+token.IdToken, 500)
		return
	}
}

func createStateNounce() string {
	nounce := make([]byte, 20)
	rand.Read(nounce)
	return "state_" + hex.EncodeToString(nounce)
}

func loginUrl(r *http.Request, exposedPort uint16, appname string) string {
	var scheme string
	if exposedPort == 8080 {
		scheme = "https"
	} else {
		scheme = "http"
	}
	return fmt.Sprintf("%s://%s/%s/api/login", scheme, r.Host, appname)
}

func createCallbackUrl(authorizationEndpoint string, clientId string, stateNounce string, callbackUrl string) (string, error) {
	fullUrl, err := url.Parse(authorizationEndpoint)
	if err != nil {
		return "", err
	}
	query := fullUrl.Query()
	query.Set("client_id", clientId)
	query.Set("state", stateNounce)
	query.Set("redirect_uri", callbackUrl)

	query.Set("session", "winssochain")
	query.Set("authIndexType", "service")
	query.Set("authIndexValue", "winssochain")

	query.Set("response_type", "code")
	query.Set("scope", "openid")

	fullUrl.RawQuery = query.Encode()

	return fullUrl.String(), nil
}

func setCookie(w http.ResponseWriter, r *http.Request, name string, value string, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		Domain:   stripPort(r.Host),
		Path:     "/",
		MaxAge:   maxAge,
		Secure:   false,
		HttpOnly: true,
	})
}

func removeCookie(w http.ResponseWriter, r *http.Request, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:    name,
		Value:   "",
		Domain:  stripPort(r.Host),
		Path:    "/",
		Expires: time.UnixMicro(0),
	})
}

func stripPort(host string) string {
	if !strings.Contains(host, ":") {
		return host
	}
	fragments := strings.Split(host, ":")
	return fragments[0]
}
