package controllers

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"webapi/app/entity"
	"webapi/app/services"
	"webapi/core"
)

// User 用户相关
type User struct {
	server  *core.Server
	service *services.Service
}

func newUser(server *core.Server, service *services.Service) *User {
	return &User{
		server:  server,
		service: service,
	}
}

// Update 更新用户信息
func (c *User) Update(ctx *fiber.Ctx) error {
	req := &entity.UserUpdateRequest{}
	if err := ctx.BodyParser(req); err != nil {
		return err
	}
	// 验证昵称格式
	if len(req.Nickname) < 2 {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "昵称长度不能小于2位")
	}
	// 验证头像地址格式
	if len(req.AvatarUrl) > 0 {
		if !strings.HasPrefix(req.AvatarUrl, "http") {
			return fiber.NewError(fiber.StatusUnprocessableEntity, "头像地址必须以http或https开头")
		}
	}
	req.Uid = ctx.UserContext().Value("uid").(int)
	resp, err := c.service.User.Update(ctx, req)
	if err != nil {
		return err
	}
	return ctx.JSON(resp)
}

// Info 获取用户信息
func (c *User) Info(ctx *fiber.Ctx) error {
	req := &entity.UserInfoRequest{}
	if err := ctx.QueryParser(req); err != nil {
		return err
	}
	if req.Uid <= 0 {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "用户id必须大于0")
	}
	resp, err := c.service.User.Info(ctx, req)
	if err != nil {
		return err
	}
	return ctx.JSON(resp)
}

// List 获取用户列表
func (c *User) List(ctx *fiber.Ctx) error {
	req := &entity.UserListRequest{}
	if err := ctx.QueryParser(req); err != nil {
		return err
	}
	if req.PageSize <= 0 {
		req.PageSize = 12
	}
	if req.OffsetId < 0 {
		req.OffsetId = 0
	}
	resp, err := c.service.User.List(ctx, req)
	if err != nil {
		return err
	}
	return ctx.JSON(resp)
}
