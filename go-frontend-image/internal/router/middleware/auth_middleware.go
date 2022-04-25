package middleware

import (
	"context"
	"errors"
	"fmt"
	"frontend-image/internal/config"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/url"
	"strings"
)

type AuthMiddlewareContext struct {
	http.Handler
	AppData *config.AppData
}

type Middleware func(handler http.Handler) http.Handler
type AuthMiddlewares struct {
	TokenExtraction Middleware
	Authentication  Middleware
}

const TOKEN_KEY = "auth_middleware_token"

func setContextToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, TOKEN_KEY, token)
}
func getContextToken(ctx context.Context) string {
	tokenStr, ok := ctx.Value(TOKEN_KEY).(string)
	if !ok {
		return ""
	}
	return tokenStr
}

func CreateAuthMiddleware(appData *config.AppData) AuthMiddlewares {
	amc := AuthMiddlewareContext{AppData: appData}
	return AuthMiddlewares{
		TokenExtraction: amc.tokenExtractionMiddleware,
		Authentication:  amc.authMiddleware,
	}
}

func (c AuthMiddlewareContext) tokenExtractionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		token, err := c.getToken(request)
		if err == nil && len(token) > 0 {
			nContext := setContextToken(request.Context(), token)
			next.ServeHTTP(writer, request.WithContext(nContext))
		} else {
			next.ServeHTTP(writer, request)
		}
	})
}

func (c AuthMiddlewareContext) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		token := getContextToken(request.Context())
		jwt, err := c.AppData.OidcClient.Verify(token)
		if err != nil {
			loginUrl := c.AppData.Config.DelegatedLoginUrl + "?url=" + url.QueryEscape(getCurrentUrl(request))
			http.Redirect(writer, request, loginUrl, http.StatusFound)
			return
		}

		log.Info().Msgf("Got token from %s", jwt.Subject())

		next.ServeHTTP(writer, request)
	})
}

const HEADER_TOKEN_LOCATION = "header"

func (c AuthMiddlewareContext) getToken(request *http.Request) (string, error) {
	tokenLocation := c.AppData.Config.AuthTokenResolver
	if tokenLocation == HEADER_TOKEN_LOCATION {
		header := request.Header.Get("Authorization")
		value := strings.TrimPrefix(header, "Bearer ")
		if len(value) == 0 {
			return "", errors.New("authorization header had zero-length value")
		}
		return value, nil
	} else {
		cookie, err := request.Cookie(tokenLocation)
		if err != nil {
			return "", nil
		}
		value := cookie.Value
		if len(value) == 0 {
			return "", errors.New("authorization cookie had zero-length value")
		}
		return value, nil
	}
}

func getCurrentUrl(request *http.Request) string {
	rHost := request.Host
	uHost := request.URL.Host
	uHostname := request.URL.Hostname()
	uSchema := request.URL.Scheme
	uPort := request.URL.Port()
	uPath := request.URL.Path
	log.Info().Msgf("%s %s %s %s %s %s", rHost, uHost, uHostname, uSchema, uPort, uPath)
	scheme := "http"
	return fmt.Sprintf("%s://%s%s", scheme, rHost, uPath)
}
