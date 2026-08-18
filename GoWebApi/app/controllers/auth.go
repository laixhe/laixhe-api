package controllers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/laixhe/gonet/utils"
	"github.com/laixhe/gonet/xfiber"

	"webapi/app/entity"
	"webapi/app/services"
	"webapi/app/util"
	"webapi/core"
	"webapi/core/middlewares"
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

// validateEmailAndPassword 验证邮箱和密码格式
func validateEmailAndPassword(email, password string) error {
	if !utils.IsEmail(email) {
		return xfiber.ParamError("邮箱格式错误")
	}
	// 密码规则 (长度 6~64 位、仅含字母 数字 _ @ $) 统一由 util.IsPassword 校验,
	// 正则已包含长度约束 (^[a-zA-Z0-9_@$]{6,64}$, 见 app/util/regexp.go), 避免规则多处维护;
	// 长度上限 64 防止超长密码被 bcrypt 静默截断 (bcrypt 只取前 72 字节)。
	if !util.IsPassword(password) {
		return xfiber.ParamError("密码格式错误，需 6~64 位，只能包含字母 数字 _ @ $")
	}
	return nil
}

// Register
// @Summary	注册
// @Tags     Auth
// @Accept   json
// @Produce  json
// @Param    req    body      entity.AuthRegisterRequest  true  "请求参数"
// @Success  200    {object}  entity.AuthRegisterResponse
// @Failure  400    {object}  core.Error
// @Failure  422    {object}  core.Error
// @Failure  500    {object}  core.Error
// @Router   /api/v1/auth/register [post]
func (c *Auth) Register(ctx fiber.Ctx) error {
	req := &entity.AuthRegisterRequest{}
	if err := ctx.Bind().WithAutoHandling().JSON(req); err != nil {
		return err
	}
	// 校验昵称与邮箱密码格式
	if err := validateNickname(req.Nickname); err != nil {
		return err
	}
	if err := validateEmailAndPassword(req.Email, req.Password); err != nil {
		return err
	}
	resp, err := c.service.Auth.Register(ctx, req)
	if err != nil {
		return err
	}
	return ctx.JSON(resp)
}

// Login
// @Summary	登录
// @Tags     Auth
// @Accept   json
// @Produce  json
// @Param    req    body      entity.AuthLoginRequest  true  "请求参数"
// @Success  200    {object}  entity.AuthLoginResponse
// @Failure  400    {object}  core.Error
// @Failure  422    {object}  core.Error
// @Failure  500    {object}  core.Error
// @Router   /api/v1/auth/login [post]
func (c *Auth) Login(ctx fiber.Ctx) error {
	req := &entity.AuthLoginRequest{}
	if err := ctx.Bind().WithAutoHandling().JSON(req); err != nil {
		return err
	}
	if err := validateEmailAndPassword(req.Email, req.Password); err != nil {
		return err
	}
	resp, err := c.service.Auth.Login(ctx, req)
	if err != nil {
		return err
	}
	return ctx.JSON(resp)
}

// Refresh
// @Summary	刷新Jwt
// @Tags     Auth
// @Accept   json
// @Produce  json
// @Param    Authorization header    string  true  "Bearer XXX令牌"
// @Success  200    {object}  entity.AuthRefreshResponse
// @Failure  400    {object}  core.Error
// @Failure  401    {object}  core.Error
// @Failure  500    {object}  core.Error
// @Router   /api/v1/auth/refresh [post]
func (c *Auth) Refresh(ctx fiber.Ctx) error {
	jwtClaims, err := middlewares.GetJwtClaims(ctx)
	if err != nil {
		return err
	}
	req := &entity.AuthRefreshRequest{Uid: jwtClaims.Uid}
	resp, err := c.service.Auth.Refresh(ctx, req)
	if err != nil {
		return err
	}
	return ctx.JSON(resp)
}
