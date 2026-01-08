package controllers

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"

	"webapi/app/entity"
	"webapi/app/services"
	"webapi/core"
	"webapi/core/middlewares"
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

// Update
// @Summary	更新用户信息
// @Tags     User
// @Accept   json
// @Produce  json
// @Param    Authorization header    string  true  "Bearer XXX令牌"
// @Param    req    body      entity.UserUpdateRequest  true  "请求参数"
// @Success  200    {object}  entity.User
// @Failure  400    {object}  core.Error
// @Failure  500    {object}  core.Error
// @Router   /api/v1/user/update [post]
func (c *User) Update(ctx fiber.Ctx) error {
	jwtClaims, err := middlewares.GetJwtClaims(ctx)
	if err != nil {
		return err
	}
	req := &entity.UserUpdateRequest{}
	if err = ctx.Bind().JSON(req); err != nil {
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
	req.Uid = jwtClaims.Uid
	resp, err := c.service.User.Update(ctx, req)
	if err != nil {
		return err
	}
	return ctx.JSON(resp)
}

// Info
// @Summary	获取用户信息
// @Tags     User
// @Accept   json
// @Produce  json
// @Param    req    body      entity.UserInfoRequest  true  "请求参数"
// @Success  200    {object}  entity.User
// @Failure  400    {object}  core.Error
// @Failure  500    {object}  core.Error
// @Router   /api/v1/user/info [get]
func (c *User) Info(ctx fiber.Ctx) error {
	req := &entity.UserInfoRequest{}
	if err := ctx.Bind().Query(req); err != nil {
		return err
	}
	log.WithContext(ctx.Context()).Debug(req)
	if req.Uid <= 0 {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "参数错误")
	}
	resp, err := c.service.User.Info(ctx, req)
	if err != nil {
		return err
	}
	return ctx.JSON(resp)
}

// List
// @Summary	获取用户列表
// @Tags     User
// @Accept   json
// @Produce  json
// @Param    req    query     entity.UserListRequest  true  "请求参数"
// @Success  200    {object}  entity.UserListResponse
// @Failure  400    {object}  core.Error
// @Failure  500    {object}  core.Error
// @Router   /api/v1/user/list [get]
func (c *User) List(ctx fiber.Ctx) error {
	req := &entity.UserListRequest{}
	if err := ctx.Bind().Query(req); err != nil {
		return err
	}
	log.WithContext(ctx.Context()).Debugf("%#v", req)
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 12
	}
	resp, err := c.service.User.List(ctx, req)
	if err != nil {
		return err
	}
	return ctx.JSON(resp)
}
