package core

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/laixhe/gonet/xlog"
)

// Error 只用于 swagger doc 生成
// Deprecated: 是 fiber.Error 的副本
type Error struct {
	Message string `json:"message" validate:"required"`
	Code    int    `json:"code" validate:"required"`
}

func (e *Error) Error() string {
	return e.Message
}

// ErrorHandler 自定义错误处理器。
//
// *fiber.Error (参数校验/鉴权/限流/超时等业务错误) 原样返回给客户端;
// 其余未知错误 (如数据库异常) 记录服务端日志后统一返回固定 500 文案,
// 避免将内部实现细节泄露给客户端。
func ErrorHandler(log *xlog.ZClient) fiber.ErrorHandler {
	return func(ctx fiber.Ctx, err error) error {
		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			return ctx.Status(fiberErr.Code).JSON(fiberErr)
		}
		log.Errorf("unhandled error request_id=%s path=%s error=%v",
			ctx.GetRespHeader(fiber.HeaderXRequestID), ctx.Path(), err)
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(fiber.NewError(fiber.StatusInternalServerError, "internal server error"))
	}
}
