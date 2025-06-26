package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

// UseErrorDefault 默认错误处理
func (m *Middleware) UseErrorDefault() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		switch errType := err.(type) {
		case *fiber.Error:
			code = errType.Code
		default:
			err = fiber.NewError(code, err.Error())
		}
		// ctx.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
		return ctx.Status(code).JSON(err)
	}
}
