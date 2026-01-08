package app

import (
	"webapi/app/controllers"
	"webapi/app/models/dao"
	"webapi/app/services"
	"webapi/core"
)

type App struct {
	Controller *controllers.Controller
	Service    *services.Service
	Dao        *dao.Dao
}

func NewApp(server *core.Server) *App {
	modelDao := dao.NewDao(server)
	service := services.NewService(server, modelDao)
	return &App{
		Controller: controllers.NewController(server, service),
		Service:    service,
		Dao:        modelDao,
	}
}
