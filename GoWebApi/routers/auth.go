package routers

import (
	"github.com/gofiber/fiber/v2"
)

// Auth 鉴权相关
func (router *Router) Auth(groupApiV1 fiber.Router) {
	groupRouter := groupApiV1.Group("auth")
	{
		groupRouter.Post("login", router.controller.Auth.Login)
	}
}
