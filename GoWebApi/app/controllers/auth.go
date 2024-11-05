package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/laixhe/gonet/ginx"
	"github.com/laixhe/gonet/logx"

	pbAuth "webapi/api/gen/auth"
	"webapi/app/services"
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
// @Param    body   body      auth.RegisterRequest   ture "请求body参数"
// @Success  200    {object}  auth.RegisterResponse
// @Router   /api/auth/register [post]
func (a *Auth) Register(c *gin.Context) {
	req := &pbAuth.RegisterRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		ginx.ErrorJSON(c, err)
		return
	}
	logx.Infof("req:%s", req)

	resp, err := a.service.AuthRegister(c, req)
	if err != nil {
		ginx.ErrorJSON(c, err)
		return
	}

	logx.Infof("resp:%s", resp)
	ginx.SuccessJSON(c, resp)
}

// Login 登录
//
// @Summary	登录用户
// @Accept   json
// @Produce  json
// @Param    body   body      auth.LoginRequest   ture "请求body参数"
// @Success  200    {object}  auth.LoginResponse
// @Router   /api/auth/login [post]
func (a *Auth) Login(c *gin.Context) {
	req := &pbAuth.LoginRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		ginx.ErrorJSON(c, err)
		return
	}
	logx.Infof("req:%s", req)

	resp, err := a.service.AuthLogin(c, req)
	if err != nil {
		ginx.ErrorJSON(c, err)
		return
	}

	logx.Infof("resp:%s", resp)
	ginx.SuccessJSON(c, resp)
}

// Refresh 刷新Jwt
//
// @Summary	刷新Jwt
// @Accept   json
// @Produce  json
// @Param Authorization header string false "Bearer token令牌"
// @Success  200    {object}  auth.RefreshResponse
// @Router   /api/auth/refresh [post]
func (a *Auth) Refresh(c *gin.Context) {
	req := &pbAuth.RefreshRequest{}
	resp, err := a.service.AuthRefresh(c, req)
	if err != nil {
		ginx.ErrorJSON(c, err)
		return
	}

	logx.Infof("resp:%s", resp)
	ginx.SuccessJSON(c, resp)
}
