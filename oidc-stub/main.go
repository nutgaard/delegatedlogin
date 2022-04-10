package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"net/http"
)

func main() {
	port := "8080"
	privateKey, publicKey, err := createKeyPair()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	jwks := createJWKS(publicKey)
	jwksJson, err := json.MarshalIndent(jwks, "", "  ")
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	appContext := AppContext{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
		Jwks:       jwks,
		JwksJson:   jwksJson,
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
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(OIDCConfig{
		Url:                   "http://oidc-stub:8080/.well-known/jwks.json",
		TokenEndpoint:         "http://oidc-stub:8080/oauth/token",
		AuthorizationEndpoint: "http://localhost:8080/authorize",
		Issuer:                "stub",
	})
	if err != nil {
		w.WriteHeader(500)
		return
	}
}

func (context AppContext) jwksHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(context.JwksJson)
	if err != nil {
		w.WriteHeader(500)
		return
	}
}

func (context AppContext) authorizationHandler(w http.ResponseWriter, r *http.Request) {
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
