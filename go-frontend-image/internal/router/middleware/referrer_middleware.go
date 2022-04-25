package middleware

import (
	"net/http"
)

func ReferrerMiddleware(referrer string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set("Referrer-Policy", referrer)
			next.ServeHTTP(writer, request)
		})
	}
}
