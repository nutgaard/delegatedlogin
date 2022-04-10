package main

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"time"
)

func createKeyPair() (jwk.RSAPrivateKey, jwk.RSAPublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	publicKey := &privateKey.PublicKey

	jwkPrivateKey, err := jwk.New(privateKey)
	if err != nil {
		return nil, nil, err
	}
	jwkPublicKey, err := jwk.New(publicKey)
	if err != nil {
		return nil, nil, err
	}

	err = jwkPublicKey.Set(jwk.KeyIDKey, uuid.New().String())
	if err != nil {
		return nil, nil, err
	}
	err = jwkPublicKey.Set(jwk.KeyUsageKey, "sig")
	if err != nil {
		return nil, nil, err
	}

	return jwkPrivateKey.(jwk.RSAPrivateKey), jwkPublicKey.(jwk.RSAPublicKey), nil
}

func createJWKS(key jwk.Key) jwk.Set {
	jwks := jwk.NewSet()
	jwks.Add(key)
	return jwks
}

func (context AppContext) createSignedJWT() string {
	token := jwt.New()
	_ = token.Set(jwt.IssuerKey, "stub")
	_ = token.Set(jwt.AudienceKey, "foo")
	_ = token.Set(jwt.SubjectKey, "Z999999")
	_ = token.Set(jwt.IssuedAtKey, time.Now())
	_ = token.Set(jwt.ExpirationKey, time.Now().Add(10*time.Minute))

	signedToken, err := jwt.Sign(token, jwa.RS256, context.PrivateKey)
	if err != nil {
		panic("Could not create token")
	}

	return string(signedToken)
}
