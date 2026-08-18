package controllers

import (
	"unicode/utf8"

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
//
// 使用 RuneCountInString 按"字符"统计而非 len() 的"字节"数,
// 否则中文等多字节字符会被误判 (如 7 个汉字=21 字节, 会被"不超过 20 位"拒绝)。
func validateNickname(nickname string) error {
	if utf8.RuneCountInString(nickname) < 2 {
		return xfiber.ParamError("昵称长度不能小于2位")
	}
	if utf8.RuneCountInString(nickname) > 20 {
		return xfiber.ParamError("昵称长度不能超过20位")
	}
	return nil
}

// normalizePagination 归一化分页参数 (与 Rust/TS/PHP 端钳制逻辑保持一致)
//
// page     <= 0 时视为 1 (默认第一页);
// pageSize <= 0 时视为 12 (默认每页条数);
// pageSize >  100 时钳制为 100 (防止超大 page_size 触发全量查询)。
func normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 12
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}
