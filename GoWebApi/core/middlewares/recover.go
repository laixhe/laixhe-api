package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// UseRecover 恢复中间件
func (m *Middleware) UseRecover(app *fiber.App) {
	app.Use(recover.New())
}
