package middleware

import (
	"net/http"
)

func CSPMiddleware(directives string, reportOnly bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if reportOnly {
				writer.Header().Set("Content-Security-Policy-Report-Only", directives)
			} else {
				writer.Header().Set("Content-Security-Policy", directives)
			}
			next.ServeHTTP(writer, request)
		})
	}
}
