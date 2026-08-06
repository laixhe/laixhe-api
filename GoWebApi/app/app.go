package app

import (
	"webapi/app/controllers"
	"webapi/app/services"
	"webapi/core"
)

// App 应用聚合，持有 Controller 和 Service 层实例
type App struct {
	Controller *controllers.Controller
	Service    *services.Service
}

// NewApp 创建应用实例，初始化 Service → Controller 依赖链
func NewApp(server *core.Server) *App {
	service := services.NewService(server)
	return &App{
		Controller: controllers.NewController(server, service),
		Service:    service,
	}
}
