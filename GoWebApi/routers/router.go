package routers

import (
	"context"

	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"webapi/app/controllers"
	"webapi/core"
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
	app := fiber.New()
	// 中间件-日志相关
	requestIdKey := router.server.Config().RequestIdKey
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: router.server.Log().Logger(),
		Fields: []string{"ip", "latency", "status", requestIdKey, "method", "url"},
	}))
	app.Use(requestid.New())
	app.Use(func(c *fiber.Ctx) error {
		ctx := context.WithValue(c.UserContext(), requestIdKey, c.GetRespHeader(fiber.HeaderXRequestID))
		c.SetUserContext(ctx)
		return c.Next()
	})
	// 中间件
	app.Use(cors.New())
	app.Use(recover.New())
	// init Server
	router.server.Init(app)
	// 路由
	groupApi := app.Group("api")
	{
		groupApiV1 := groupApi.Group("v1")
		{
			router.Auth(groupApiV1) // 鉴权相关
		}
	}
	return router
}

func (router *Router) Listen() error {
	return router.server.Listen()
}
