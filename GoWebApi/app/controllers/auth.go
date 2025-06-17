package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"

	"webapi/app/entity"
	"webapi/app/services"
	"webapi/core"
)

// Auth 鉴权相关
type Auth struct {
	server  *core.Server
	service *services.Service
}

func newAuth(server *core.Server, service *services.Service) *Auth {
	return &Auth{
		server:  server,
		service: service,
	}
}

// Login 登录
func (a *Auth) Login(c *fiber.Ctx) error {
	req := &entity.AuthLoginRequest{}
	if err := c.BodyParser(req); err != nil {
		return err
	}
	resp, err := a.service.Auth.Login(c, req)
	if err != nil {
		return c.SendString(err.Error())
	}
	log.WithContext(c.UserContext()).Warnw("T------- api v1 auth login -------", "resp", resp)
	return c.JSON(resp)
}
