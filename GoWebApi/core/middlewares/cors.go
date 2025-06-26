package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// UseCors 跨域中间件
func (m *Middleware) UseCors(app *fiber.App) {
	app.Use(cors.New())
}
