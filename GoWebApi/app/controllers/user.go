package controllers

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/laixhe/gonet/xfiber"

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
// @Param    req    body      entity.UserUpdateRequest  true  "请求参数（Uid 由 JWT 提供）"
// @Success  200    {object}  entity.User
// @Failure  401    {object}  core.Error
// @Failure  400    {object}  core.Error
// @Failure  422    {object}  core.Error
// @Failure  500    {object}  core.Error
// @Router   /api/v1/user/update [post]
func (c *User) Update(ctx fiber.Ctx) error {
	jwtClaims, err := middlewares.GetJwtClaims(ctx)
	if err != nil {
		return err
	}
	req := &entity.UserUpdateRequest{}
	if err = ctx.Bind().WithAutoHandling().JSON(req); err != nil {
		return err
	}
	// 校验昵称格式
	if err := validateNickname(req.Nickname); err != nil {
		return err
	}
	// 验证头像地址格式
	if len(req.AvatarUrl) > 255 {
		return xfiber.ParamError("头像地址长度不能超过255位")
	}
	// 必须精确以 http:// 或 https:// 开头 (不用 HasPrefix("http"), 否则 httpxxx:// 也能通过)
	if len(req.AvatarUrl) > 0 {
		if !strings.HasPrefix(req.AvatarUrl, "http://") && !strings.HasPrefix(req.AvatarUrl, "https://") {
			return xfiber.ParamError("头像地址必须以http或https开头")
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
// @Param    req    query     entity.UserInfoRequest  true  "请求参数"
// @Success  200    {object}  entity.User
// @Failure  400    {object}  core.Error
// @Failure  422    {object}  core.Error
// @Failure  500    {object}  core.Error
// @Router   /api/v1/user/info [get]
func (c *User) Info(ctx fiber.Ctx) error {
	req := &entity.UserInfoRequest{}
	if err := ctx.Bind().WithAutoHandling().Query(req); err != nil {
		return err
	}
	if req.Uid <= 0 {
		return xfiber.ParamError("无效的用户ID")
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
// @Failure  422    {object}  core.Error
// @Failure  500    {object}  core.Error
// @Router   /api/v1/user/list [get]
func (c *User) List(ctx fiber.Ctx) error {
	req := &entity.UserListRequest{}
	if err := ctx.Bind().WithAutoHandling().Query(req); err != nil {
		return err
	}
	// 归一化分页参数: page<=0→1, page_size<=0→12, page_size>100→100 (见 normalizePagination)
	req.Page, req.PageSize = normalizePagination(req.Page, req.PageSize)
	resp, err := c.service.User.List(ctx, req)
	if err != nil {
		return err
	}
	return ctx.JSON(resp)
}
