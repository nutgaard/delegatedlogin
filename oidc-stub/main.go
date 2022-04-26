package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/rs/zerolog/log"
	"net/http"
)

func main() {
	var err error
	port := "8080"
	jwks := createJWKS()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	appContext := AppContext{
		Jwks:   jwks,
		Config: loadConfig(),
	}
	log.Info().Msg(appContext.createSignedJWT())

	http.HandleFunc("/.well-known/openid-configuration", appContext.oidcHandler)
	http.HandleFunc("/.well-known/jwks.json", appContext.jwksHandler)
	http.HandleFunc("/authorize", appContext.authorizationHandler)
	http.HandleFunc("/oauth/token", appContext.tokenHandler)

	log.Info().Msg("Listening to " + port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
}

func (context AppContext) oidcHandler(w http.ResponseWriter, _ *http.Request) {
	log.Info().Msg("200 GET /.well-known/openid-configuration")
	w.Header().Set("Content-Type", "application/json")
	domain := "localhost"
	if context.Config.DockerCompose {
		domain = "oidc-stub"
	}
	err := json.NewEncoder(w).Encode(OIDCConfig{
		Url:                   fmt.Sprintf("http://%s:8080/.well-known/jwks.json", domain),
		TokenEndpoint:         fmt.Sprintf("http://%s:8080/oauth/token", domain),
		AuthorizationEndpoint: "http://localhost:8080/authorize",
		Issuer:                "stub",
	})
	if err != nil {
		w.WriteHeader(500)
		return
	}
}

func (context AppContext) jwksHandler(w http.ResponseWriter, _ *http.Request) {
	log.Info().Msg("200 GET /.well-known/jwks.json")
	w.Header().Set("Content-Type", "application/json")
	publicKeys, err := jwk.PublicSetOf(context.Jwks)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	jwksJson, _ := json.MarshalIndent(publicKeys, "", "  ")
	_, err = w.Write(jwksJson)
	if err != nil {
		w.WriteHeader(500)
		return
	}
}

func (context AppContext) authorizationHandler(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("200 GET /authorize")
	redirectUri := r.URL.Query().Get("redirect_uri")
	state := r.URL.Query().Get("state")
	if len(redirectUri) == 0 {
		http.Error(w, "redirect_url missing", http.StatusBadRequest)
	} else if len(state) == 0 {
		http.Error(w, "state missing", http.StatusBadRequest)
	} else {
		code := uuid.New().String()
		http.Redirect(w, r,
			fmt.Sprintf("%s?code=%s&state=%s", redirectUri, code, state),
			http.StatusFound,
		)
	}
}

func (context AppContext) tokenHandler(w http.ResponseWriter, _ *http.Request) {
	log.Info().Msg("200 GET /oauth/token")
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(TokenExchangeResult{
		IdToken:      context.createSignedJWT(),
		AccessToken:  "random access",
		RefreshToken: "refresh_token",
	})
	if err != nil {
		w.WriteHeader(500)
		return
	}
}
