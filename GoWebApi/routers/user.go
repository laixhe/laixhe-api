package routers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/laixhe/gonet/xfiber"
)

// User 用户相关
func (r *Router) User(routerApi fiber.Router) {
	groupRouter := routerApi.Group("user")
	{
		groupRouter.Get("info", r.app.Controller.User.Info) // 获取用户信息
		groupRouter.Get("list", r.app.Controller.User.List) // 获取用户列表
	}
	groupRouter.Use(xfiber.UseJwt(r.middleware.UseJwtConfig))
	{
		groupRouter.Post("update", r.app.Controller.User.Update) // 更新用户信息
	}
}
