package app

import (
	"webapi/app/controllers"
	"webapi/app/services"
	"webapi/core"
)

type App struct {
	Controller *controllers.Controller
	Service    *services.Service
}

func NewApp(server *core.Server) *App {
	service := services.NewService(server)
	return &App{
		Controller: controllers.NewController(server, service),
		Service:    service,
	}
}
