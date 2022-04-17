package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type logFields struct {
	CallId string
	Error  error
}

func (fields *logFields) MarshalZerologObject(e *zerolog.Event) {
	e.Str("callId", fields.CallId)

	if fields.Error != nil {
		e.Err(fields.Error)
	}
}

type MaskingConfig struct {
	Pattern     string
	Replacement string
}

func ZerologMiddleware(config MaskingConfig) func(ctx *fiber.Ctx) error {
	pattern := regexp.MustCompile(config.Pattern)
	replacement := config.Replacement

	return func(ctx *fiber.Ctx) error {
		start := time.Now()

		callId := ctx.Locals(CALL_ID).(string)

		fields := &logFields{
			CallId: callId,
		}

		chainError := ctx.Next()

		fields.Error = chainError
		statusCode := ctx.Response().StatusCode()

		var logmsg strings.Builder

		_, _ = logmsg.WriteString(strconv.Itoa(statusCode))
		_, _ = logmsg.WriteString(" ")
		_, _ = logmsg.WriteString(http.StatusText(statusCode))
		_, _ = logmsg.WriteString(" ")
		_, _ = logmsg.WriteString(fmt.Sprintf("%v", time.Since(start).Milliseconds()))
		_, _ = logmsg.WriteString("ms")
		_, _ = logmsg.WriteString(" ")
		_, _ = logmsg.WriteString(ctx.Method())
		_, _ = logmsg.WriteString(" ")
		_, _ = logmsg.WriteString(ctx.Path())

		log.
			Info().
			EmbedObject(fields).
			Msg(mask(logmsg.String(), pattern, replacement))

		return chainError
	}
}

func mask(txt string, pattern *regexp.Regexp, replacement string) string {
	return pattern.ReplaceAllStringFunc(txt, func(match string) string {
		return strings.Repeat(replacement, len(match))
	})
}
