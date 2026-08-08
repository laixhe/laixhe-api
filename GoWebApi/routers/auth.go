package routers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/laixhe/gonet/xfiber"
)

// Auth 鉴权相关
func (r *Router) Auth(routerApi fiber.Router) {
	groupRouter := routerApi.Group("auth")
	{
		// 公开接口 (不受 JWT 保护)
		groupRouter.Post("register", r.app.Controller.Auth.Register) // 注册
		groupRouter.Post("login", r.app.Controller.Auth.Login)       // 登录
	}
	// Use(Jwt) 只作用于其后注册的路由: 仅 refresh 需要 JWT
	groupRouter.Use(xfiber.UseJwt(r.middleware.UseJwtConfig))
	{
		groupRouter.Post("refresh", r.app.Controller.Auth.Refresh) // 刷新 Jwt
	}
}
