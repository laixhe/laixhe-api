package routers

import (
	"github.com/gofiber/fiber/v2"

	"webapi/core/middlewares"
)

// User 用户相关
func (router *Router) User(groupApiV1 fiber.Router, secretKey string) {
	groupRouter := groupApiV1.Group("user")
	{
		groupRouter.Get("info", router.controller.User.Info) // 获取用户信息
		groupRouter.Get("list", router.controller.User.List) // 获取用户列表
	}
	groupRouter.Use(middlewares.UseJwt(secretKey), middlewares.UseJwtClaims())
	{
		groupRouter.Post("update", router.controller.User.Update) // 更新用户信息
	}
}
