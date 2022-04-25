package main

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/rs/zerolog/log"
	"time"
)

func createJWKS() jwk.Set {
	jwks := jwk.NewSet()
	privateKey, _ := createRSAKey()
	jwks.Add(privateKey)
	return jwks
}

func createRSAKey() (jwk.RSAPrivateKey, error) {
	rsaPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	privateKey := jwk.NewRSAPrivateKey()
	err = privateKey.Set(jwk.KeyIDKey, uuid.New().String())
	err = privateKey.FromRaw(rsaPrivateKey)
	return privateKey, err
}

func (context AppContext) createSignedJWT() string {
	token := jwt.New()
	_ = token.Set(jwt.IssuerKey, "stub")
	_ = token.Set(jwt.AudienceKey, "foo")
	_ = token.Set(jwt.SubjectKey, "Z999999")
	_ = token.Set(jwt.IssuedAtKey, time.Now())
	_ = token.Set(jwt.ExpirationKey, time.Now().Add(10*time.Minute))
	key, _ := context.Jwks.Get(0)

	signedToken, err := jwt.Sign(token, jwa.RS256, key)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not create token")
	}

	return string(signedToken)
}
