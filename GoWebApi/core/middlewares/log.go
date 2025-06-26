package middlewares

import (
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// UseLog 日志中间件
func (m *Middleware) UseLog(app *fiber.App, zapLogger *zap.Logger) {
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: zapLogger,
		Fields: []string{"ip", "latency", "status", m.RequestIdKey, "method", "url"},
	}))
}
