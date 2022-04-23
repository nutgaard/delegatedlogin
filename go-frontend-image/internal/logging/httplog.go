package logging

import (
	"github.com/go-chi/httplog"
	"github.com/rs/zerolog"
)

func NewHttpLogger(appName string) zerolog.Logger {
	opts := httplog.Options{
		Concise:  true,
		LogLevel: "info",
		JSON:     true,
	}
	return httplog.NewLogger(appName, opts)
}
