package services

import (
	"webapi/app/models"
	"webapi/core"

	"github.com/gofiber/fiber/v2"
)

type Auth struct {
	server *core.Server
	model  *models.Model
}

func NewAuth(server *core.Server) *Auth {
	return &Auth{
		server: server,
		model:  models.NewModel(server),
	}
}

func (ss *Auth) Login(c *fiber.Ctx) (*models.User, error) {
	return ss.model.GetUserByNickname(c.UserContext(), "laixhe")
}
