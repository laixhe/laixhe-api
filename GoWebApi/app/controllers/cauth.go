package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"

	"webapi/app/services"
	"webapi/core"
)

type Auth struct {
	auth *services.Auth
}

func NewAuth(server *core.Server) *Auth {
	return &Auth{
		auth: services.NewAuth(server),
	}
}

// Login 登录
func (cs *Auth) Login(c *fiber.Ctx) error {
	// return errors.New("T---------------------XX")
	user, err := cs.auth.Login(c)
	if err != nil {
		return c.SendString(err.Error())
	}
	log.WithContext(c.UserContext()).Warnw("T------- api v1 auth login -------", "user", user)
	return c.SendString(c.BaseURL())
}
