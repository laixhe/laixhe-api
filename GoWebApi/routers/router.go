package routers

import (
	"github.com/gofiber/fiber/v3"

	"webapi/app"
	"webapi/core"
	"webapi/core/middlewares"
	"webapi/docs"
)

// Router 业务路由
type Router struct {
	server     *core.Server
	app        *app.App
	middleware *middlewares.Middleware
}

// NewRouter 创建路由实例并初始化所有路由
func NewRouter(server *core.Server) *Router {
	r := &Router{
		server:     server,
		app:        app.NewApp(server),
		middleware: middlewares.NewMiddleware(server.Config().Jwt.SecretKey),
	}
	return r.init()
}

// init 注册 API 路由组及 Swagger 文档端点
func (r *Router) init() *Router {
	// 路由
	groupApi := r.server.Server().App().Group("api")
	{
		groupApiV1 := groupApi.Group("v1")
		{
			groupApiV1.Get("swagger.json", func(ctx fiber.Ctx) error {
				return ctx.SendString(docs.JsonSwagger)
			})
			groupApiV1.Get("swagger.yaml", func(ctx fiber.Ctx) error {
				return ctx.SendString(docs.YamlSwagger)
			})
			r.Auth(groupApiV1) // 鉴权相关
			r.User(groupApiV1) // 用户相关
		}
	}
	return r
}

// HttpStart 启动Http服务
func (r *Router) HttpStart() error {
	return r.server.Server().Listen(r.server.Config().Http.Addr())
}
