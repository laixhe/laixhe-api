package routers

import (
	"github.com/gofiber/fiber/v2"

	"webapi/app/controllers"
	"webapi/core"
	"webapi/core/middlewares"
)

// Router 业务路由
type Router struct {
	server     *core.Server
	controller *controllers.Controller
}

func NewRouter(server *core.Server) *Router {
	router := &Router{
		server:     server,
		controller: controllers.NewController(server),
	}
	return router.init()
}

func (router *Router) init() *Router {
	router.server.HttpMiddleware(func(app *fiber.App) {
		middlewares.UseLog(app, router.server.Log().Logger(), core.RequestIdKey)
		middlewares.UseRequestId(app, core.RequestIdKey)
		middlewares.UseCors(app)
		middlewares.UseRecover(app)
	})
	// 路由
	groupApi := router.server.HttpGroup("api")
	{
		groupApiV1 := groupApi.Group("v1")
		{
			router.Auth(groupApiV1, router.server.Config().Jwt.SecretKey) // 鉴权相关
			router.User(groupApiV1, router.server.Config().Jwt.SecretKey) // 用户相关
		}
	}
	return router
}

// HttpStart 启动Http服务
func (router *Router) HttpStart() error {
	return router.server.HttpStart()
}
