package middlewares

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

// UseRequestId 请求ID中间件
func UseRequestId(app *fiber.App, requestIdKey string) {
	app.Use(requestid.New())
	app.Use(func(ctx *fiber.Ctx) error {
		newCtx := context.WithValue(ctx.UserContext(),
			requestIdKey,
			ctx.GetRespHeader(fiber.HeaderXRequestID))
		ctx.SetUserContext(newCtx)
		return ctx.Next()
	})
}
