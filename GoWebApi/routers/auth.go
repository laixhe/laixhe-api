package routers

import (
	"github.com/gofiber/fiber/v2"

	"webapi/core/middlewares"
)

// Auth 鉴权相关
func (router *Router) Auth(groupApiV1 fiber.Router, secretKey string) {
	groupRouter := groupApiV1.Group("auth")
	{
		groupRouter.Post("register", router.controller.Auth.Register) // 注册
		groupRouter.Post("login", router.controller.Auth.Login)       // 登录
	}
	groupRouter.Use(middlewares.UseJwt(secretKey), middlewares.UseJwtClaims())
	{
		groupRouter.Post("refresh", router.controller.Auth.Refresh) // 刷新Jwt
	}
}
