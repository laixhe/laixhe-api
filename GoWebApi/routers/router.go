package routers

import (
	"context"

	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"webapi/core"
)

// Router 路由
func Router(server *core.Server) *fiber.App {
	router := fiber.New()
	// 中间件-日志相关
	requestidKey := server.Config().RequestIdKey
	router.Use(fiberzap.New(fiberzap.Config{
		Logger: server.Log().Logger(),
		Fields: []string{"ip", "latency", "status", requestidKey, "method", "url"},
	}))
	router.Use(requestid.New())
	router.Use(func(c *fiber.Ctx) error {
		ctx := context.WithValue(c.UserContext(), requestidKey, c.GetRespHeader(fiber.HeaderXRequestID))
		c.SetUserContext(ctx)
		return c.Next()
	})
	// 中间件
	router.Use(cors.New())
	router.Use(recover.New())
	// init Server
	server.SetApp(router)
	server.Init()
	// 路由
	api := router.Group("api")
	{
		apiV1 := api.Group("v1")
		{
			AuthRouter(server, apiV1)
		}
	}
	return router
}
