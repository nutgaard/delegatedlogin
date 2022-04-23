package router

import (
	"net/http"
)

func (handler *Handler) Authorization(next http.Handler) http.Handler {
	fn := func(writer http.ResponseWriter, request *http.Request) {

		next.ServeHTTP(writer, request)
	}
	return http.HandlerFunc(fn)
}
