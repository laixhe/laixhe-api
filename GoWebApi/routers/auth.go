package routers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/laixhe/gonet/xfiber"
)

// Auth 鉴权相关
func (r *Router) Auth(routerApi fiber.Router) {
	groupRouter := routerApi.Group("auth")
	{
		groupRouter.Post("register", r.app.Controller.Auth.Register) // 注册
		groupRouter.Post("login", r.app.Controller.Auth.Login)       // 登录
	}
	groupRouter.Use(xfiber.UseJwt(r.middleware.UseJwtConfig))
	{
		groupRouter.Post("refresh", r.app.Controller.Auth.Refresh) // 刷新 Jwt
	}
}
