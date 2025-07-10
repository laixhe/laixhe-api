package middlewares

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

// UseErrorDefault 默认错误处理
func UseErrorDefault() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		var errType *fiber.Error
		switch {
		case errors.As(err, &errType):
			code = errType.Code
		default:
			err = fiber.NewError(code, err.Error())
		}
		return ctx.Status(code).JSON(err)
	}
}
