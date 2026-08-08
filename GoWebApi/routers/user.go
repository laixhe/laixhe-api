package routers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/laixhe/gonet/xfiber"
)

// User 用户相关
func (r *Router) User(routerApi fiber.Router) {
	groupRouter := routerApi.Group("user")
	{
		// 公开用户接口 (不受 JWT 保护)
		groupRouter.Get("info", r.app.Controller.User.Info) // 获取用户信息
		groupRouter.Get("list", r.app.Controller.User.List) // 获取用户列表
	}
	// Use(Jwt) 只作用于其后注册的路由: 仅 update 需要 JWT
	// xfiber.UseJwt 底层为 gofiber/contrib jwt: 校验令牌签名与过期时间
	groupRouter.Use(xfiber.UseJwt(r.middleware.UseJwtConfig))
	{
		groupRouter.Post("update", r.app.Controller.User.Update) // 更新用户信息
	}
}
