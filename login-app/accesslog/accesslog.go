package accesslog

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

type logginResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *logginResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func Decorate(handler func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lrw := &logginResponseWriter{w, http.StatusOK}
		handler(lrw, r)

		url := r.URL.String()
		method := r.Method
		status := lrw.statusCode

		elapsed := time.Since(start)
		log.Info().
			Int64("time_ms", elapsed.Milliseconds()).
			Msg(fmt.Sprintf("%s %d %s", method, status, url))
	}
}
