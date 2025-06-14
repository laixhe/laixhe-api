package routers

import (
	"github.com/gofiber/fiber/v2"

	"webapi/app/controllers"
	"webapi/core"
)

// AuthRouter 鉴权相关
func AuthRouter(server *core.Server, router fiber.Router) {
	groupRouter := router.Group("auth")
	cs := controllers.NewAuth(server)
	{
		groupRouter.Post("login", cs.Login)
	}
}
