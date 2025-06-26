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
	middleware *middlewares.Middleware
	controller *controllers.Controller
}

func NewRouter(server *core.Server) *Router {
	router := &Router{
		server:     server,
		middleware: middlewares.NewMiddleware(server.Config().RequestIdKey, server.Config().Jwt.SecretKey),
		controller: controllers.NewController(server),
	}
	return router.init()
}

func (router *Router) init() *Router {
	app := fiber.New(fiber.Config{
		ErrorHandler: router.middleware.UseErrorDefault(),
	})
	router.server.Init(app)
	// 中间件
	router.middleware.UseLog(app, router.server.Log().Logger())
	router.middleware.UseRequestId(app)
	router.middleware.UseCors(app)
	router.middleware.UseRecover(app)
	// 路由
	groupApi := app.Group("api")
	{
		groupApiV1 := groupApi.Group("v1")
		{
			router.Auth(groupApiV1) // 鉴权相关
			router.User(groupApiV1) // 用户相关
		}
	}
	return router
}

func (router *Router) Listen() error {
	return router.server.Listen()
}
