package middleware

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

const CALL_ID = "callId"

func CallIdMiddleware(next http.Handler) http.Handler {
	fn := func(writer http.ResponseWriter, request *http.Request) {
		value := uuid.New().String()

		ctx := request.Context()
		context.WithValue(ctx, CALL_ID, value)
		next.ServeHTTP(writer, request.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}
