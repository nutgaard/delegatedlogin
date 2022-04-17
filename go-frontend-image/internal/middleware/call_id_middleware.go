package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const CALL_ID = "callId"

func CallIdMiddleware() func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		value := ctx.Get(CALL_ID)
		if value == "" {
			value = uuid.New().String()
			ctx.Locals(CALL_ID, value)
		}

		return ctx.Next()
	}
}
