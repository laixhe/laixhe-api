package controllers

import (
	"webapi/app/services"
	"webapi/core"
)

// Controller 业务控制器
type Controller struct {
	Auth *Auth
	User *User
}

func NewController(server *core.Server, service *services.Service) *Controller {
	return &Controller{
		Auth: newAuth(server, service),
		User: newUser(server, service),
	}
}
