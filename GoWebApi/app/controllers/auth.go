package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/laixhe/gonet/xgin"
	"github.com/laixhe/gonet/xlog"

	"webapi/app/services"
	"webapi/core"
	"webapi/protocol/gen/pbauth"
)

type Auth struct {
	service *services.Service
}

func NewAuth() *Auth {
	return &Auth{
		service: services.NewService(),
	}
}

// Register 注册
//
// @Summary	注册用户
// @Accept   json
// @Produce  json
// @Param    body   body      pbauth.RegisterRequest   ture "请求body参数"
// @Success  200    {object}  pbauth.RegisterResponse
// @Router   /api/auth/register [post]
func (a *Auth) Register(c *gin.Context) {
	req := &pbauth.RegisterRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		core.JSONErrorParse(c, err)
		return
	}
	//
	xlog.Info(fmt.Sprintf("req:%v", req), xgin.ZapField(c)...)
	//
	resp, err := a.service.AuthRegister(c, req)
	if err != nil {
		core.JSONError(c, err)
		return
	}
	//
	xlog.Info(fmt.Sprintf("resp:%v", resp), xgin.ZapField(c)...)
	//
	xgin.SuccessJSON(c, resp)
}

// Login 登录
//
// @Summary	登录用户
// @Accept   json
// @Produce  json
// @Param    body   body      pbauth.LoginRequest   ture "请求body参数"
// @Success  200    {object}  pbauth.LoginResponse
// @Router   /api/auth/login [post]
func (a *Auth) Login(c *gin.Context) {
	req := &pbauth.LoginRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		core.JSONErrorParse(c, err)
		return
	}
	//
	xlog.Info(fmt.Sprintf("req:%v", req), xgin.ZapField(c)...)
	//
	resp, err := a.service.AuthLogin(c, req)
	if err != nil {
		core.JSONError(c, err)
		return
	}
	//
	xlog.Info(fmt.Sprintf("resp:%v", resp), xgin.ZapField(c)...)
	//
	xgin.SuccessJSON(c, resp)
}

// Refresh 刷新Jwt
//
// @Summary	刷新Jwt
// @Accept   json
// @Produce  json
// @Param Authorization header string false "Bearer token令牌"
// @Success  200    {object}  pbauth.RefreshResponse
// @Router   /api/auth/refresh [post]
func (a *Auth) Refresh(c *gin.Context) {
	req := &pbauth.RefreshRequest{}
	resp, err := a.service.AuthRefresh(c, req)
	if err != nil {
		core.JSONError(c, err)
		return
	}
	//
	xlog.Info(fmt.Sprintf("resp:%v", resp), xgin.ZapField(c)...)
	//
	xgin.SuccessJSON(c, resp)
}
