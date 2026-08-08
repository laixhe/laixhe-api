package routers

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/timeout"

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
		middleware: middlewares.NewMiddleware(server),
	}
	return r.init()
}

// init 注册 API 路由组及 Swagger 文档端点
//
// 中间件执行顺序 (从外到内, 由 xfiber.New 与 init 中的注册顺序决定):
// requestId → 访问日志 → panic 恢复 → CORS → gzip 压缩 → 请求超时(408) → IP 限流(429) → 业务路由
func (r *Router) init() *Router {
	// 请求超时中间件 (超过 http.timeout 秒未完成返回 408)
	r.server.Fiber().App().Use(timeout.New(func(c fiber.Ctx) error {
		return c.Next()
	}, timeout.Config{
		Timeout: time.Duration(r.server.Config().Http.Timeout) * time.Second,
		// 超时响应统一 JSON (避免纯文本)
		OnTimeout: func(c fiber.Ctx) error {
			return c.Status(fiber.StatusRequestTimeout).
				JSON(fiber.NewError(fiber.StatusRequestTimeout, "Request Timeout"))
		},
	}))
	// 全局限流中间件 (基于客户端 IP, 超过阈值返回 429 统一 JSON)
	r.server.Fiber().App().Use(r.middleware.RateLimit)
	// 路由
	groupApi := r.server.Fiber().App().Group("api")
	{
		groupApiV1 := groupApi.Group("v1")
		{
			groupApiV1.Get("swagger.json", func(ctx fiber.Ctx) error {
				// 显式设置 Content-Type 与缓存头
				ctx.Set(fiber.HeaderContentType, "application/json")
				ctx.Set(fiber.HeaderCacheControl, "public, max-age=300")
				return ctx.SendString(docs.JsonSwagger)
			})
			groupApiV1.Get("swagger.yaml", func(ctx fiber.Ctx) error {
				// 显式设置 Content-Type 与缓存头
				ctx.Set(fiber.HeaderContentType, "application/x-yaml")
				ctx.Set(fiber.HeaderCacheControl, "public, max-age=300")
				return ctx.SendString(docs.YamlSwagger)
			})
			// 健康检查 (含数据库探测, 限流已豁免该路径)
			groupApiV1.Get("health", r.app.Controller.Health.Health)
			r.Auth(groupApiV1) // 鉴权相关
			r.User(groupApiV1) // 用户相关
		}
	}
	return r
}

// HttpStart 启动Http服务 (支持 Ctrl+C / SIGTERM 优雅停机)
func (r *Router) HttpStart() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return r.server.Fiber().App().Listen(r.server.Config().Http.Address(), fiber.ListenConfig{
		GracefulContext: ctx,
	})
}
