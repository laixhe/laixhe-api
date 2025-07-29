package controllers

import (
	"github.com/duke-git/lancet/v2/validator"
	"github.com/gofiber/fiber/v2"

	"webapi/app/entity"
	"webapi/app/services"
	"webapi/app/utils"
	"webapi/core"
)

// Auth 鉴权相关
type Auth struct {
	server  *core.Server
	service *services.Service
}

func newAuth(server *core.Server, service *services.Service) *Auth {
	return &Auth{
		server:  server,
		service: service,
	}
}

// @Summary	注册
// @Tags     Auth
// @Accept   json
// @Produce  json
// @Param    req       body     entity.AuthRegisterRequest  true  "请求参数"
// @Success  200       {object}  entity.AuthRegisterResponse
// @Router   /api/v1/auth/register [post]
func (c *Auth) Register(ctx *fiber.Ctx) error {
	req := &entity.AuthRegisterRequest{}
	if err := ctx.BodyParser(req); err != nil {
		return err
	}
	// 验证昵称格式
	if len(req.Nickname) < 2 {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "昵称长度不能小于2位")
	}
	// 验证邮箱格式
	if !validator.IsEmail(req.Email) {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "邮箱格式错误")
	}
	if len(req.Password) < 6 {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "密码长度不能小于6位")
	}
	// 验证密码格式
	if !utils.IsPassword(req.Password) {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "密码格式错误，只能包含字母 数字 _ @ $")
	}
	resp, err := c.service.Auth.Register(ctx, req)
	if err != nil {
		return err
	}
	return ctx.JSON(resp)
}

// @Summary	登录
// @Tags     Auth
// @Accept   json
// @Produce  json
// @Param    req       body     entity.AuthLoginRequest  true  "请求参数"
// @Success  200       {object}  entity.AuthLoginResponse
// @Router   /api/v1/auth/login [post]
func (c *Auth) Login(ctx *fiber.Ctx) error {
	req := &entity.AuthLoginRequest{}
	if err := ctx.BodyParser(req); err != nil {
		return err
	}
	// 验证邮箱格式
	if !validator.IsEmail(req.Email) {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "邮箱格式错误")
	}
	if len(req.Password) < 6 {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "密码长度不能小于6位")
	}
	// 验证密码格式
	if !utils.IsPassword(req.Password) {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "密码格式错误，只能包含字母 数字 _ @ $")
	}
	resp, err := c.service.Auth.Login(ctx, req)
	if err != nil {
		return err
	}
	return ctx.JSON(resp)
}

// @Summary	刷新Jwt
// @Tags     Auth
// @Accept   json
// @Produce  json
// @Param Authorization header string true "Bearer token令牌"
// @Success  200       {object}  entity.AuthRefreshResponse
// @Router   /api/v1/auth/refresh [post]
func (c *Auth) Refresh(ctx *fiber.Ctx) error {
	uid := ctx.UserContext().Value("uid").(int)
	req := &entity.AuthRefreshRequest{Uid: uid}
	resp, err := c.service.Auth.Refresh(ctx, req)
	if err != nil {
		return err
	}
	return ctx.JSON(resp)
}
