package controllers

import (
	"github.com/laixhe/gonet/xfiber"

	"webapi/app/services"
	"webapi/core"
)

// Controller 业务控制器
type Controller struct {
	Auth   *Auth
	User   *User
	Health *Health
}

// NewController 创建控制器实例，初始化 Auth、User 和 Health 子控制器
func NewController(server *core.Server, service *services.Service) *Controller {
	return &Controller{
		Auth:   newAuth(server, service),
		User:   newUser(server, service),
		Health: newHealth(server),
	}
}

// validateNickname 校验昵称长度 (注册与更新用户信息共用)
func validateNickname(nickname string) error {
	if len(nickname) < 2 {
		return xfiber.ParamError("昵称长度不能小于2位")
	}
	if len(nickname) > 20 {
		return xfiber.ParamError("昵称长度不能超过20位")
	}
	return nil
}
